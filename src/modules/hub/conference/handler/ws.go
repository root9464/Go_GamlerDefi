package conference_ws

import (
	"encoding/json"

	"github.com/gofiber/contrib/socketio"
	"github.com/pion/webrtc/v4"
	conference_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/util"
)

func (h *WSHandler) Disconect(ep *socketio.EventPayload) {
	h.conference_usecase.Disconect(ep)
}

func (h *WSHandler) Connect(ep *socketio.EventPayload) {
	h.logger.Infof("New connection: socket_id=%s, request_id=%s", ep.Kws.GetUUID(), conference_utils.GenerateRequestID())
}

func (h *WSHandler) SetupSocketEventHandlers(ep *socketio.EventPayload) {
	requestID := conference_utils.GenerateRequestID()
	message := &conference_utils.WebsocketMessage{}
	if err := json.Unmarshal(ep.Data, &message); err != nil {
		h.logger.Errorf("Failed to unmarshal message, error: %v, request_id: %s", err, requestID)
		return
	}

	roomID := ep.Kws.GetStringAttribute("room_id")
	if roomID == "" {
		h.logger.Errorf("No room ID associated with connection, socket_id: %s, request_id: %s", ep.Kws.GetUUID(), requestID)
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
		h.logger.Errorf("No PeerConnection found for socket, socket_id: %s, request_id: %s", ep.Kws.GetUUID(), requestID)
		return
	}

	switch message.Event {
	case "offer":
		offer := webrtc.SessionDescription{}
		if err := json.Unmarshal([]byte(message.Data), &offer); err != nil {
			h.logger.Errorf("Failed to unmarshal offer, error: %v, request_id: %s", err, requestID)
			return
		}
		if err := conn.Pc.SetRemoteDescription(offer); err != nil {
			h.logger.Errorf("Failed to set remote description, error: %v, request_id: %s", err, requestID)
			return
		}
		answer, err := conn.Pc.CreateAnswer(nil)
		if err != nil {
			h.logger.Errorf("Failed to create answer, error: %v, request_id: %s", err, requestID)
			return
		}
		if err = conn.Pc.SetLocalDescription(answer); err != nil {
			h.logger.Errorf("Failed to set local description, error: %v, request_id: %s", err, requestID)
			return
		}
		answerString, err := json.Marshal(answer)
		if err != nil {
			h.logger.Errorf("Failed to marshal answer, error: %v, request_id: %s", err, requestID)
			return
		}
		if err := conference_utils.WriteJSON(conn.Kws, &conn.Lock, &conference_utils.WebsocketMessage{
			Event: "answer",
			Data:  string(answerString),
		}); err != nil {
			h.logger.Errorf("Failed to send answer, error: %v, request_id: %s", err, requestID)
			return
		}

		h.logger.Infof("Offer processed and answer sent, socket_id: %s, request_id: %s", ep.Kws.GetUUID(), requestID)

	case "candidate":
		candidate := webrtc.ICECandidateInit{}
		if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
			h.logger.Errorf("Failed to unmarshal candidate, error: %v, request_id: %s", err, requestID)
			return
		}
		if err := conn.Pc.AddICECandidate(candidate); err != nil {
			h.logger.Errorf("Failed to add ICE candidate, error: %v, request_id: %s", err, requestID)
			return
		}

	case "answer":
		answer := webrtc.SessionDescription{}
		if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
			h.logger.Errorf("Failed to unmarshal answer, error: %v, request_id: %s", err, requestID)
			return
		}
		if conn.Pc.SignalingState() != webrtc.SignalingStateHaveLocalOffer {
			h.logger.Warnf("Skipping SetRemoteDescription due to invalid signaling state, state: %s, socket_id: %s, request_id: %s", conn.Pc.SignalingState().String(), conn.Kws.GetUUID(), requestID)
			return
		}
		if err := conn.Pc.SetRemoteDescription(answer); err != nil {
			h.logger.Errorf("Failed to set remote description, error: %v, request_id: %s", err, requestID)
			return
		}

		h.logger.Infof("Answer processed, socket_id: %s, request_id: %s", conn.Kws.GetUUID(), requestID)
	default:
		h.logger.Warnf("Unknown message event: %s", message.Event)
	}
}

func (h *WSHandler) ConferenceHandler(kws *socketio.Websocket) {
	requestID := conference_utils.GenerateRequestID()
	roomID := kws.Params("session_id")
	if roomID == "" {
		h.logger.Errorf("Room ID not provided, socket_id: %s, request_id: %s", kws.GetUUID(), requestID)
		return
	}

	kws.SetAttribute("room_id", roomID)

	mediaEngine := &webrtc.MediaEngine{}
	if err := mediaEngine.RegisterDefaultCodecs(); err != nil {
		h.logger.Errorf("Failed to register default codecs, error: %v, request_id: %s", err, requestID)
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
		h.logger.Errorf("Failed to create PeerConnection, error: %v, request_id: %s", err, requestID)
		return
	}

	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := pc.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			h.logger.Errorf("Failed to add transceiver, error: %v, request_id: %s", err, requestID)
			return
		}
	}

	conn := h.conference_usecase.CreateConnection(roomID, pc, kws)
	room := h.conference_usecase.GetOrCreateRoom(roomID, requestID, conn)
	h.conference_usecase.SetubWebRTC(conn, room, requestID)

	h.conference_usecase.SignalPeerConnections(requestID, roomID)

}
