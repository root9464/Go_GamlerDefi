package conference_ws

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/socketio"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (h *WSHandler) socketErr(ep *socketio.EventPayload, err error) {
	errByte := []byte(err.Error())
	ep.Kws.Emit(errByte)
}

func (h *WSHandler) ConfirenceSocketHandler(kws *socketio.Websocket) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}},
	}
	mediaEngine := webrtc.MediaEngine{}
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:    webrtc.MimeTypeOpus,
			ClockRate:   48000,
			Channels:    2,
			SDPFmtpLine: "minptime=10;useinbandfec=1",
		},
		PayloadType: 111,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		kws.Emit([]byte(fmt.Sprintf("register audio codec: %v", err)), socketio.TextMessage)
	}
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypeVP8,
			ClockRate: 90000,
		},
		PayloadType: 96,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		kws.Emit([]byte(fmt.Sprintf("register video codec: %v", err)), socketio.TextMessage)
	}
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	pc, err := api.NewPeerConnection(config)
	if err != nil {
		kws.Emit([]byte(fmt.Sprintf("create peer connection: %v", err)), socketio.TextMessage)
	}
	defer pc.Close()

	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := pc.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			kws.Emit([]byte(fmt.Sprintf("add transceiver: %v", err)), socketio.TextMessage)
		}
	}

	writer := &conference_utils.ThreadSafeWriter{Conn: kws}
	peer := &conference_utils.PeerConnection{PC: pc, Conn: writer}
	h.conference_usecase.AddPeer(pc, writer)

	pc.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}
		candidateData, err := json.Marshal(i.ToJSON())
		if err != nil {
			h.logger.Errorf("marshal candidate: %v", err)
			return
		}
		if err := writer.WriteJSON(conference_utils.WebsocketMessage{
			Event: "candidate",
			Data:  string(candidateData),
		}); err != nil {
			h.logger.Errorf("send candidate: %v", err)
		}
	})

	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		if state == webrtc.PeerConnectionStateClosed {
			err := h.conference_usecase.SignalPeers()
			if err != nil {
				h.logger.Errorf("error transfer signals in peer %s", err.Error())
				return
			}
		}
	})

	pc.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		trackLocal, err := h.conference_usecase.AddTrack(peer, t)
		if err != nil {
			h.logger.Errorf("add track: %v", err)
			return
		}
		defer h.conference_usecase.RemoveTrack(trackLocal)

		bufferSize := 1500
		if t.Kind() == webrtc.RTPCodecTypeAudio {
			bufferSize = 500
		}
		buf := make([]byte, bufferSize)
		for {
			n, _, err := t.Read(buf)
			if err != nil {
				return
			}
			var pkt rtp.Packet
			if err := pkt.Unmarshal(buf[:n]); err != nil {
				h.logger.Errorf("unmarshal RTP packet: %v", err)
				return
			}
			if t.Kind() != webrtc.RTPCodecTypeAudio {
				pkt.Extension = false
				pkt.Extensions = nil
			}
			if err := trackLocal.WriteRTP(&pkt); err != nil {
				return
			}
		}
	})

	if err := h.conference_usecase.SignalPeers(); err != nil {
		h.logger.Errorf("signal peers: %v", err)
	}

	socketio.On("candidate", func(ep *socketio.EventPayload) {
		var candidate webrtc.ICECandidateInit
		if err := json.Unmarshal([]byte(ep.Data), &candidate); err != nil {
			h.logger.Errorf("unmarshal candidate: %v", err)
			ep.Kws.Emit([]byte(fmt.Sprintf("unmarshal candidate: %v", err)), socketio.TextMessage)
		}
		if err := pc.AddICECandidate(candidate); err != nil {
			h.logger.Errorf("add ICE candidate: %v", err)
			ep.Kws.Emit([]byte(fmt.Sprintf("add ICE candidate: %v", err)), socketio.TextMessage)
		}
	})

	socketio.On("answer", func(ep *socketio.EventPayload) {
		var answer webrtc.SessionDescription
		if err := json.Unmarshal([]byte(ep.Data), &answer); err != nil {
			h.logger.Errorf("unmarshal answer: %v", err)
			ep.Kws.Emit([]byte(fmt.Sprintf("unmarshal answer: %v", err)), socketio.TextMessage)
		}
		if err := pc.SetRemoteDescription(answer); err != nil {
			h.logger.Errorf("set remote description: %v", err)
			ep.Kws.Emit([]byte(fmt.Sprintf("set remote description: %v", err)), socketio.TextMessage)
		}
	})
}
