package conference_ws

import (
	"encoding/json"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	"github.com/pion/webrtc/v4"
	hub_entity "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (h *WSHandler) ConferenceWebsocketHandler(c *fiber.Ctx) error {
	socketio.On("connect", func(ep *socketio.EventPayload) {
		h.logger.Info("New connection",
			"socket_id", ep.Kws.GetUUID(),
			"request_id", conference_utils.GenerateRequestID())
	})

	socketio.On("disconnect", func(ep *socketio.EventPayload) {
		h.conference_usecase.Disconect(ep)
	})

	socketio.On("message", func(ep *socketio.EventPayload) {
		requestID := conference_utils.GenerateRequestID()
		message := &conference_utils.WebsocketMessage{}
		if err := json.Unmarshal(ep.Data, &message); err != nil {
			h.logger.Error("Failed to unmarshal message",
				"error", err,
				"request_id", requestID)
			return
		}

		roomID := ep.Kws.GetStringAttribute("room_id")
		if roomID == "" {
			h.logger.Error("No room ID associated with connection",
				"uuid", ep.Kws.GetUUID(),
				"request_id", requestID)
			return
		}

		room := h.conference_usecase.GetOrCreateRoom(roomID, requestID, nil)
		room.Lock.RLock()
		var conn *hub_entity.Connection
		for _, c := range room.Connections {
			if c.Kws.GetUUID() == ep.Kws.GetUUID() {
				conn = c
				break
			}
		}
		room.Lock.RUnlock()

		if conn == nil || conn.Pc == nil {
			h.logger.Error("No PeerConnection found",
				"uuid", ep.Kws.GetUUID(),
				"request_id", requestID)
			return
		}

		switch message.Event {
		case "offer":
			offer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &offer); err != nil {
				h.logger.Error("Failed to unmarshal offer",
					"error", err,
					"request_id", requestID)
				return
			}
			if err := conn.Pc.SetRemoteDescription(offer); err != nil {
				h.logger.Error("Failed to set remote description",
					"error", err,
					"request_id", requestID)
				return
			}
			answer, err := conn.Pc.CreateAnswer(nil)
			if err != nil {
				h.logger.Error("Failed to create answer",
					"error", err,
					"request_id", requestID)
				return
			}
			if err = conn.Pc.SetLocalDescription(answer); err != nil {
				h.logger.Error("Failed to set local description",
					"error", err,
					"request_id", requestID)
				return
			}
			answerString, err := json.Marshal(answer)
			if err != nil {
				h.logger.Error("Failed to marshal answer",
					"error", err,
					"request_id", requestID)
				return
			}
			if err := conference_utils.WriteJSON(conn.Kws, &conn.Lock, &conference_utils.WebsocketMessage{
				Event: "answer",
				Data:  string(answerString),
			}); err != nil {
				h.logger.Error("Failed to send answer",
					"error", err,
					"request_id", requestID)
				return
			}

			h.logger.Info("Offer processed and answer sent",
				"uuid", ep.Kws.GetUUID(),
				"request_id", requestID)

		case "candidate":
			candidate := webrtc.ICECandidateInit{}
			if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
				h.logger.Error("Failed to unmarshal candidate",
					"error", err,
					"request_id", requestID)
				return
			}
			if err := conn.Pc.AddICECandidate(candidate); err != nil {
				h.logger.Error("Failed to add ICE candidate",
					"error", err,
					"request_id", requestID)
				return
			}

		case "answer":
			answer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
				h.logger.Error("Failed to unmarshal answer",
					"error", err,
					"request_id", requestID)
				return
			}
			if conn.Pc.SignalingState() != webrtc.SignalingStateHaveLocalOffer {
				h.logger.Warn("Skipping SetRemoteDescription due to invalid signaling state",
					"state", conn.Pc.SignalingState().String(),
					"uuid", conn.Kws.GetUUID(),
					"request_id", requestID)
				return
			}
			if err := conn.Pc.SetRemoteDescription(answer); err != nil {
				h.logger.Error("Failed to set remote description",
					"error", err,
					"request_id", requestID)
				return
			}

			h.logger.Info("Answer processed",
				"uuid", conn.Kws.GetUUID(),
				"request_id", requestID)
			h.conference_usecase.SignalPeerConnections(requestID, conn.RoomID)

		default:
			h.logger.Warn("Unknown message event",
				"event", message.Event,
				"request_id", requestID)
		}

	})
	return socketio.New(func(kws *socketio.Websocket) {
		// if atomic.LoadUint32(&serverRunning) == 0 {
		// 	h.logger.Info("Connection attempt while server is shutting down",
		// 		"uuid", kws.GetUUID())
		// 	return
		// }

		requestID := conference_utils.GenerateRequestID()
		roomID := kws.Query("room_id")
		if roomID == "" {
			h.logger.Error("Room ID not provided",
				"uuid", kws.GetUUID(),
				"request_id", requestID)
			return
		}

		kws.SetAttribute("room_id", roomID)

		mediaEngine := &webrtc.MediaEngine{}
		if err := mediaEngine.RegisterDefaultCodecs(); err != nil {
			h.logger.Error("Failed to register default codecs",
				"request_id", requestID)
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
			h.logger.Error("Failed to create PeerConnection",
				"error", err,
				"request_id", requestID)
			return
		}

		for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
			if _, err := pc.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionRecvonly,
			}); err != nil {
				h.logger.Error("Failed to add transceiver",
					"error", err,
					"request_id", requestID)
				return
			}
		}

		conn := h.conference_usecase.CreateConnection(roomID, pc, kws)
		room := h.conference_usecase.GetOrCreateRoom(roomID, requestID, conn)
		h.conference_usecase.SetubWebRTC(conn, room, requestID)

		h.conference_usecase.SignalPeerConnections(requestID, roomID)
	})(c)
}
