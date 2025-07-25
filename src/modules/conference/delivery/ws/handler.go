package conference_ws_handler

import (
	"encoding/json"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/conference/utils"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type PeerService interface {
	AddPeer(pc *webrtc.PeerConnection, conn *conference_utils.ThreadSafeWriter)
	SignalPeers() error
}

type TrackService interface {
	AddTrack(peer *conference_utils.PeerConnection, t *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error)
	RemoveTrack(track *webrtc.TrackLocalStaticRTP)
}

type WSHandler struct {
	logger        *logger.Logger
	peer_service  PeerService
	track_service TrackService
}

func NewWSHanler(
	logger *logger.Logger,
	peerService PeerService,
	trackService TrackService,
) *WSHandler {
	return &WSHandler{
		logger:        logger,
		peer_service:  peerService,
		track_service: trackService,
	}
}

func (s *WSHandler) ConferenceWebsocketHandler(c *fiber.Ctx) error {
	return socketio.New(func(conn *socketio.Websocket) {
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
			s.logger.Errorf("register audio codec: %v", err)
			return
		}

		if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
			RTPCodecCapability: webrtc.RTPCodecCapability{
				MimeType:  webrtc.MimeTypeVP8,
				ClockRate: 90000,
			},
			PayloadType: 96,
		}, webrtc.RTPCodecTypeVideo); err != nil {
			s.logger.Errorf("register video codec: %v", err)
			return
		}

		api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
		pc, err := api.NewPeerConnection(config)
		if err != nil {
			s.logger.Errorf("create peer connection: %v", err)
			return
		}
		defer pc.Close()

		for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
			if _, err := pc.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionRecvonly,
			}); err != nil {
				s.logger.Errorf("add transceiver: %v", err)
				return
			}
		}

		writer := &conference_utils.ThreadSafeWriter{Conn: conn.Conn}
		peer := &conference_utils.PeerConnection{PC: pc, Conn: writer}
		s.peer_service.AddPeer(pc, writer)

		pc.OnICECandidate(func(i *webrtc.ICECandidate) {
			if i == nil {
				return
			}
			candidateData, err := json.Marshal(i.ToJSON())
			if err != nil {
				s.logger.Errorf("marshal candidate: %v", err)
				return
			}
			if err := writer.WriteJSON(conference_utils.WebsocketMessage{
				Event: "candidate",
				Data:  string(candidateData),
			}); err != nil {
				s.logger.Errorf("send candidate: %v", err)
			}
		})

		pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
			if state == webrtc.PeerConnectionStateClosed {
				if err := s.peer_service.SignalPeers(); err != nil {
					s.logger.Errorf("error transfer signals in peer %s", err.Error())
				}
			}
		})

		pc.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
			trackLocal, err := s.track_service.AddTrack(peer, t)
			if err != nil {
				s.logger.Errorf("add track: %v", err)
				return
			}
			defer s.track_service.RemoveTrack(trackLocal)

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
					s.logger.Errorf("unmarshal RTP packet: %v", err)
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

		if err := s.peer_service.SignalPeers(); err != nil {
			s.logger.Errorf("signal peers: %v", err)
		}

		for {
			_, data, err := conn.Conn.ReadMessage()
			if err != nil {
				return
			}
			var msg conference_utils.WebsocketMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				s.logger.Errorf("unmarshal message: %v", err)
				continue
			}
			switch msg.Event {
			case "candidate":
				var candidate webrtc.ICECandidateInit
				if err := json.Unmarshal([]byte(msg.Data), &candidate); err != nil {
					s.logger.Errorf("unmarshal candidate: %v", err)
					continue
				}
				if err := pc.AddICECandidate(candidate); err != nil {
					s.logger.Errorf("add ICE candidate: %v", err)
				}
			case "answer":
				var answer webrtc.SessionDescription
				if err := json.Unmarshal([]byte(msg.Data), &answer); err != nil {
					s.logger.Errorf("unmarshal answer: %v", err)
					continue
				}
				if err := pc.SetRemoteDescription(answer); err != nil {
					s.logger.Errorf("set remote description: %v", err)
				}
			default:
				s.logger.Errorf("unknown message event: %s", msg.Event)
			}
		}
	})(c)
}
