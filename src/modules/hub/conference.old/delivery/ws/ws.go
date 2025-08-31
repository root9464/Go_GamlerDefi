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
		case "offer":
			offer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &offer); err != nil {
				h.logger.Errorf("failed to unmarshal offer, error: %v, request id: %s", err, requestID)
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
			if err := conference_utils.WriteJSON(conn.Kws, &conn.Lock, &conference_utils.WebsocketMessage{
				Event: "answer",
				Data:  string(answerString),
			}); err != nil {
				h.logger.Errorf("failed to send answer, error: %v, request id: %s", err, requestID)
				return
			}

			h.logger.Infof("offer processed and answer sent, uuid: %s, request id: %s", ep.Kws.GetUUID(), requestID)
		case "candidate":
			candidate := webrtc.ICECandidateInit{}
			if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
				h.logger.Errorf("failed to unmarshal candidate, error: %v, request id: %s", err, requestID)
				return
			}
			if err := conn.Pc.AddICECandidate(candidate); err != nil {
				h.logger.Errorf("failed to add ICE candidate, error: %v, request id: %s", err, requestID)
				return
			}

		case "answer":
			answer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
				h.logger.Errorf("failed to unmarshal answer, error: %v, request id: %s", err, requestID)
				return
			}
			switch conn.Pc.SignalingState() {
			case webrtc.SignalingStateHaveLocalOffer, webrtc.SignalingStateHaveRemotePranswer:
				return
			case webrtc.SignalingStateStable:
				currentRemote := conn.Pc.RemoteDescription()
				if currentRemote != nil && currentRemote.SDP == answer.SDP {
					h.logger.Debug("Duplicate answer ignored",
						"state", conn.Pc.SignalingState().String(),
						"uuid", conn.Kws.GetUUID(),
						"request_id", requestID)
					return
				}
			default:
				h.logger.Warnf("skipping SetRemoteDescription due to invalid signaling state, state: %s, uuid: %s, request id: %s", conn.Pc.SignalingState().String(), conn.Kws.GetUUID(), requestID)
				return
			}

			if err := conn.Pc.SetRemoteDescription(answer); err != nil {
				h.logger.Errorf("failed to set remote description, error: %v, request id: %s", err, requestID)
				return
			}

			h.logger.Infof("answer processed, uuid: %s, request id: %s", conn.Kws.GetUUID(), requestID)
			h.conference_usecase.SignalPeerConnections(requestID, conn.RoomID)

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

		h.conference_usecase.SignalPeerConnections(requestID, roomID)
	})(c)
}
