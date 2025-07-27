package conference_ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

const (
	EventCandidate = "candidate"
	EventAnswer    = "answer"
)

type PeerSession struct {
	WPC    *webrtc.PeerConnection
	PC     *conference_utils.PeerConnection
	Conn   *socketio.Websocket
	ctx    context.Context
	cancel context.CancelFunc
}

func (ps *PeerSession) Close() {
	ps.cancel()
	ps.WPC.Close()
}

func (s *WSHandler) ConferenceWebsocketHandler(c *fiber.Ctx) error {
	hubID := c.Query("hubID")
	userID := c.Query("userID")

	if hubID == "" || userID == "" {
		s.logger.Warn("Missing hubID or userID in query")
		return errors.New("hubID and userID are required")
	}

	return socketio.New(func(conn *socketio.Websocket) {
		s.logger.Infof("New WebSocket connection for userID: %s, hubID: %s", userID, hubID)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		session, err := s.initPeerSession(ctx, conn, userID, hubID)
		if err != nil {
			s.logger.Errorf("Init peer session failed: %v", err)
			return
		}
		s.logger.Infof("Peer session initialized for userID: %s", userID)

		if err := s.conference_usecase.JoinHub(session.PC); err != nil {
			s.logger.Errorf("JoinHub failed: %v", err)
			return
		}
		s.logger.Infof("User %s joined hub %s", userID, hubID)

		defer func() {
			s.logger.Infof("Cleaning up user %s", userID)
			// if err := s.conference_usecase.LeaveHub(session.PC); err != nil {
			// 	s.logger.Errorf("LeaveHub error: %v", err)
			// }
			session.Close()
		}()

		s.conference_usecase.AddPeer(session.WPC, session.PC)

		session.WPC.OnICECandidate(func(cand *webrtc.ICECandidate) {
			if cand == nil {
				return
			}
			s.logger.Debugf("New ICE candidate gathered for user %s", userID)
			if err := s.sendCandidate(session, cand); err != nil {
				s.logger.Errorf("Send candidate error: %v", err)
			}
		})

		session.WPC.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
			s.logger.Infof("Connection state changed to %s for user %s", state.String(), userID)
			if state == webrtc.PeerConnectionStateClosed {
				s.logger.Warnf("Peer connection closed for user %s", userID)
				if err := s.conference_usecase.SignalPeers(session.PC); err != nil {
					s.logger.Errorf("SignalPeers after close error: %v", err)
				}
			}
		})

		session.WPC.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
			s.logger.Infof("Received track of type %s from user %s", t.Kind().String(), userID)
			s.handleTrack(session, t)
		})

		s.logger.Infof("Sending initial signaling for user %s", userID)
		if err := s.conference_usecase.SignalPeers(session.PC); err != nil {
			s.logger.Errorf("Initial SignalPeers error: %v", err)
		}

		for {
			msg, err := session.ReadMessage()
			if err != nil {
				s.logger.Infof("WebSocket read error or connection closed for user %s: %v", userID, err)
				return
			}

			s.logger.Debugf("Received WS message: %+v", msg)

			if err := s.handleWebsocketMessage(session, msg); err != nil {
				s.logger.Errorf("handleWebsocketMessage error: %v", err)
			}
		}
	})(c)
}

func (ps *PeerSession) ReadMessage() (conference_utils.WebsocketMessage, error) {
	type result struct {
		msg []byte
		err error
	}
	ch := make(chan result, 1)

	go func() {
		_, msg, err := ps.Conn.Conn.ReadMessage()
		ch <- result{msg, err}
	}()

	select {
	case <-ps.ctx.Done():
		return conference_utils.WebsocketMessage{}, errors.New("context canceled")
	case res := <-ch:
		if res.err != nil {
			return conference_utils.WebsocketMessage{}, res.err
		}
		var wsMsg conference_utils.WebsocketMessage
		if err := json.Unmarshal(res.msg, &wsMsg); err != nil {
			return conference_utils.WebsocketMessage{}, err
		}
		return wsMsg, nil
	}
}

func (s *WSHandler) initPeerSession(ctx context.Context, conn *socketio.Websocket, userID, hubID string) (*PeerSession, error) {
	s.logger.Infof("Initializing PeerConnection for user: %s, hub: %s", userID, hubID)

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}},
	}
	mediaEngine := webrtc.MediaEngine{}

	// Register Opus
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:    webrtc.MimeTypeOpus,
			ClockRate:   48000,
			Channels:    2,
			SDPFmtpLine: "minptime=10;useinbandfec=1",
		},
		PayloadType: 111,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		return nil, fmt.Errorf("register opus codec: %w", err)
	}

	// Register VP8
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypeVP8,
			ClockRate: 90000,
		},
		PayloadType: 96,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		return nil, fmt.Errorf("register vp8 codec: %w", err)
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))

	wpc, err := api.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("create PeerConnection: %w", err)
	}

	for _, kind := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := wpc.AddTransceiverFromKind(kind, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionSendrecv,
		}); err != nil {
			wpc.Close()
			return nil, fmt.Errorf("add transceiver (%s): %w", kind.String(), err)
		}
	}

	writer := &conference_utils.ThreadSafeWriter{Conn: conn.Conn}
	pc := &conference_utils.PeerConnection{
		PC:     wpc,
		Writer: writer,
		UserID: userID,
		HubID:  hubID,
	}

	// wpc.OnNegotiationNeeded(func() {
	// 	s.logger.Infof("Negotiation needed for user: %s", userID)
	// 	go func() {
	// 		if err := s.conference_usecase.SignalPeers(pc); err != nil {
	// 			s.logger.Errorf("Negotiation signal error: %v", err)
	// 		}
	// 	}()
	// })

	session := &PeerSession{
		WPC:  wpc,
		PC:   pc,
		Conn: conn,
	}
	session.ctx, session.cancel = context.WithCancel(ctx)

	return session, nil
}

func (s *WSHandler) sendCandidate(session *PeerSession, candidate *webrtc.ICECandidate) error {
	data, err := json.Marshal(candidate.ToJSON())
	if err != nil {
		return fmt.Errorf("marshal candidate: %w", err)
	}
	msg := conference_utils.WebsocketMessage{
		Event: EventCandidate,
		Data:  string(data),
	}
	s.logger.Debugf("Sending ICE candidate to client: %s", data)
	return session.PC.Writer.WriteJSON(msg)
}

func (s *WSHandler) handleTrack(session *PeerSession, t *webrtc.TrackRemote) {
	s.logger.Infof("HandleTrack started for track kind: %s", t.Kind().String())
	trackLocal, err := s.conference_usecase.AddTrack(session.PC, t)
	if err != nil {
		s.logger.Errorf("AddTrack failed: %v", err)
		return
	}
	defer s.conference_usecase.RemoveTrack(trackLocal, session.PC)

	bufferSize := 1500
	if t.Kind() == webrtc.RTPCodecTypeAudio {
		bufferSize = 500
	}
	buf := make([]byte, bufferSize)

	for {
		n, _, err := t.Read(buf)
		if err != nil {
			s.logger.Infof("Track read finished: %v", err)
			return
		}

		var pkt rtp.Packet
		if err := pkt.Unmarshal(buf[:n]); err != nil {
			s.logger.Errorf("Unmarshal RTP packet error: %v", err)
			continue
		}

		if t.Kind() != webrtc.RTPCodecTypeAudio {
			pkt.Extension = false
			pkt.Extensions = nil
		}

		if err := trackLocal.WriteRTP(&pkt); err != nil {
			s.logger.Errorf("Write RTP error: %v", err)
			return
		}
	}
}

func (s *WSHandler) handleWebsocketMessage(session *PeerSession, msg conference_utils.WebsocketMessage) error {
	switch msg.Event {
	case EventCandidate:
		s.logger.Infof("Received ICE candidate from client: %s", msg.Data)
		var candidate webrtc.ICECandidateInit
		if err := json.Unmarshal([]byte(msg.Data), &candidate); err != nil {
			return fmt.Errorf("unmarshal candidate: %w", err)
		}
		return session.WPC.AddICECandidate(candidate)

	case EventAnswer:
		s.logger.Infof("Received answer from client")
		var answer webrtc.SessionDescription
		if err := json.Unmarshal([]byte(msg.Data), &answer); err != nil {
			return fmt.Errorf("unmarshal answer: %w", err)
		}

		return session.WPC.SetRemoteDescription(answer)

	default:
		s.logger.Warnf("Unknown event received: %s", msg.Event)
		return nil
	}
}
