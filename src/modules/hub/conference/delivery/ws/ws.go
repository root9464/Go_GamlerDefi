// conference_ws/ws.go
package conference_ws

import (
	"encoding/json"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	"github.com/pion/webrtc/v4"
	conference_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/utils"
)

func (h *WSHandler) ConferenceWebsocketEvents() {
	socketio.On("connect", func(ep *socketio.EventPayload) {
		h.logger.Infof("new connection socket id: %s, request id: %s", ep.Kws.GetUUID(), conference_utils.GenerateRequestID())
	})

	socketio.On("disconnect", func(ep *socketio.EventPayload) {
		h.conference_usecase.Disconect(ep)
	})

	socketio.On("message", func(ep *socketio.EventPayload) {
		requestID := conference_utils.GenerateRequestID()
		message := &conference_utils.WebsocketMessage{}
		if err := json.Unmarshal(ep.Data, &message); err != nil {
			h.logger.Errorf("failed to unmarshal message, error: %v, request id: %s", err, requestID)
			return
		}

		roomID := ep.Kws.GetStringAttribute("room_id")
		if roomID == "" {
			h.logger.Errorf("no room ID associated with connection, uuid: %s, request id: %s", ep.Kws.GetUUID(), requestID)
			return
		}

		room := h.conference_usecase.GetOrCreateRoom(roomID, requestID, nil)
		room.Lock.RLock()
		var conn *conference_entity.Connection
		for _, c := range room.Connections {
			if c.Kws.GetUUID() == ep.Kws.GetUUID() {
				conn = c
				break
			}
		}
		room.Lock.RUnlock()

		if conn == nil || conn.Pc == nil {
			h.logger.Errorf("no PeerConnection found, uuid: %s, request id: %s", ep.Kws.GetUUID(), requestID)
			return
		}

		switch message.Event {
		case "request_offer":
			h.logger.Infof("request_offer received, uuid: %s, request id: %s", conn.Kws.GetUUID(), requestID)
			h.conference_usecase.SignalPeerConnections(requestID, roomID)

		case "offer":
			offer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &offer); err != nil {
				h.logger.Errorf("failed to unmarshal offer, error: %v, request id: %s", err, requestID)
				return
			}

			if conn.Pc.SignalingState() != webrtc.SignalingStateStable {
				h.logger.Warnf("cannot set remote offer in current signaling state: %s, uuid: %s, request id: %s",
					conn.Pc.SignalingState().String(), conn.Kws.GetUUID(), requestID)
				return
			}

			if err := conn.Pc.SetRemoteDescription(offer); err != nil {
				h.logger.Errorf("failed to set remote description, error: %v, request id: %s", err, requestID)
				return
			}

			answer, err := conn.Pc.CreateAnswer(nil)
			if err != nil {
				h.logger.Errorf("failed to create answer, error: %v, request id: %s", err, requestID)
				return
			}

			if err = conn.Pc.SetLocalDescription(answer); err != nil {
				h.logger.Errorf("failed to set local description, error: %v, request id: %s", err, requestID)
				return
			}

			answerString, err := json.Marshal(answer)
			if err != nil {
				h.logger.Errorf("failed to marshal answer, error: %v, request id: %s", err, requestID)
				return
			}

			if err := conference_utils.WriteJSON(conn.Kws, &conn.Mu, &conference_utils.WebsocketMessage{
				Event: "answer",
				Data:  string(answerString),
			}); err != nil {
				h.logger.Errorf("failed to send answer, error: %v, request id: %s", err, requestID)
				return
			}

			h.logger.Infof("offer processed and answer sent, uuid: %s, request id: %s", ep.Kws.GetUUID(), requestID)

		case "answer":
			answer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
				h.logger.Errorf("failed to unmarshal answer, error: %v, request id: %s", err, requestID)
				return
			}

			switch conn.Pc.SignalingState() {
			case webrtc.SignalingStateStable:
				currentRemote := conn.Pc.RemoteDescription()
				if currentRemote != nil && currentRemote.SDP == answer.SDP {
					h.logger.Debugf("duplicate answer ignored, state: %s, uuid: %s, request id: %s",
						conn.Pc.SignalingState().String(), conn.Kws.GetUUID(), requestID)
					return
				}
				h.logger.Warnf("unexpected answer in stable state, uuid: %s, request id: %s",
					conn.Kws.GetUUID(), requestID)
				return

			case webrtc.SignalingStateHaveLocalOffer:
				if err := conn.Pc.SetRemoteDescription(answer); err != nil {
					h.logger.Errorf("failed to set remote description, error: %v, request id: %s", err, requestID)
					return
				}
				h.logger.Infof("answer processed, uuid: %s, request id: %s", conn.Kws.GetUUID(), requestID)

			default:
				h.logger.Warnf("cannot set remote answer in current signaling state: %s, uuid: %s, request id: %s",
					conn.Pc.SignalingState().String(), conn.Kws.GetUUID(), requestID)
				return
			}

		case "candidate":
			candidate := webrtc.ICECandidateInit{}
			if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
				h.logger.Errorf("failed to unmarshal candidate, error: %v, request id: %s", err, requestID)
				return
			}

			conn.Mu.Lock()
			if conn.Pc.RemoteDescription() == nil {
				conn.PendingCandidates = append(conn.PendingCandidates, candidate)
				h.logger.Debugf("caching ICE candidate - no remote description, uuid: %s, request id: %s",
					conn.Kws.GetUUID(), requestID)
			} else {
				if err := conn.Pc.AddICECandidate(candidate); err != nil {
					h.logger.Errorf("failed to add ICE candidate, error: %v, request id: %s", err, requestID)
				}
			}
			conn.Mu.Unlock()

		default:
			h.logger.Warnf("unknown message event, event: %s, request id: %s", message.Event, requestID)
		}
	})
}

func (h *WSHandler) ConferenceWebsocketHandler(c *fiber.Ctx) error {
	h.ConferenceWebsocketEvents()
	return socketio.New(func(kws *socketio.Websocket) {
		requestID := conference_utils.GenerateRequestID()
		roomID := kws.Query("room_id")
		if roomID == "" {
			h.logger.Errorf("room ID not provided, uuid: %s, request id: %s", kws.GetUUID(), requestID)
			return
		}

		kws.SetAttribute("room_id", roomID)

		mediaEngine := &webrtc.MediaEngine{}
		if err := mediaEngine.RegisterDefaultCodecs(); err != nil {
			h.logger.Errorf("failed to register default codecs, request id: %s", requestID)
			return
		}

		settingEngine := webrtc.SettingEngine{}
		api := webrtc.NewAPI(
			webrtc.WithMediaEngine(mediaEngine),
			webrtc.WithSettingEngine(settingEngine),
		)

		pc, err := api.NewPeerConnection(webrtc.Configuration{
			SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
		})
		if err != nil {
			h.logger.Errorf("failed to create peer connection, error: %v, request id: %s", err, requestID)
			return
		}

		for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
			if _, err := pc.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionRecvonly,
			}); err != nil {
				h.logger.Errorf("failed to add transceiver, error: %v, request id: %s", err, requestID)
				return
			}
		}

		conn := h.conference_usecase.CreateConnection(roomID, pc, kws)
		room := h.conference_usecase.GetOrCreateRoom(roomID, requestID, conn)
		h.conference_usecase.SetubWebRTC(conn, room, requestID)
		h.logger.Infof("peer connection created, waiting for offer request, uuid: %s, request id: %s", kws.GetUUID(), requestID)
	})(c)
}
