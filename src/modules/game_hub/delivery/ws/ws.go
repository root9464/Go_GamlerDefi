package conference_ws

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

// Константы для websocket-событий
const (
	EventCandidate = "candidate"
	EventAnswer    = "answer"
)

type PeerSession struct {
	PC     *webrtc.PeerConnection
	Writer *conference_utils.ThreadSafeWriter
	Conn   *socketio.Websocket
	ctx    context.Context
	cancel context.CancelFunc
}

func (ps *PeerSession) Close() {
	ps.cancel()
	ps.PC.Close()
}

// Чтение сообщений с контекстом отмены
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

func (s *WSHandler) ConferenceWebsocketHandler(c *fiber.Ctx) error {
	return socketio.New(func(conn *socketio.Websocket) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Инициализация PeerConnection и сессии
		session, err := s.initPeerSession(ctx, conn)
		if err != nil {
			s.logger.Errorf("init peer session: %v", err)
			return
		}
		defer session.Close()

		// Регистрируем пира в конференции
		s.conference_usecase.AddPeer(session.PC, session.Writer)

		// Обработка ICE кандидатов
		session.PC.OnICECandidate(func(cand *webrtc.ICECandidate) {
			if cand == nil {
				return
			}
			if err := s.sendCandidate(session, cand); err != nil {
				s.logger.Errorf("send candidate: %v", err)
			}
		})

		// Обработка изменения состояния соединения
		session.PC.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
			if state == webrtc.PeerConnectionStateClosed {
				if err := s.conference_usecase.SignalPeers(); err != nil {
					s.logger.Errorf("signal peers after close: %v", err)
				}
			}
		})

		// Обработка входящих медиатреков
		session.PC.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
			s.handleTrack(session, t)
		})

		// Начальная сигнализация пиров
		if err := s.conference_usecase.SignalPeers(); err != nil {
			s.logger.Errorf("initial signal peers: %v", err)
		}

		// Главный цикл чтения WebSocket-сообщений от клиента
		for {
			msg, err := session.ReadMessage()
			if err != nil {
				s.logger.Infof("read message error or connection closed: %v", err)
				return
			}

			if err := s.handleWebsocketMessage(session, msg); err != nil {
				s.logger.Errorf("handle websocket message: %v", err)
			}
		}
	})(c)
}

// Инициализация PeerConnection и сессии с трансиверами и кодеками
func (s *WSHandler) initPeerSession(ctx context.Context, conn *socketio.Websocket) (*PeerSession, error) {
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
		return nil, err
	}

	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypeVP8,
			ClockRate: 90000,
		},
		PayloadType: 96,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		return nil, err
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))

	pc, err := api.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	// Добавляем трансиверы для приема аудио и видео
	for _, kind := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := pc.AddTransceiverFromKind(kind, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			pc.Close()
			return nil, err
		}
	}

	writer := &conference_utils.ThreadSafeWriter{Conn: conn.Conn}
	peerSession := &PeerSession{
		PC:     pc,
		Writer: writer,
		Conn:   conn,
	}
	peerSession.ctx, peerSession.cancel = context.WithCancel(ctx)

	return peerSession, nil
}

// Отправка ICE-кандидата клиенту через WebSocket
func (s *WSHandler) sendCandidate(session *PeerSession, candidate *webrtc.ICECandidate) error {
	data, err := json.Marshal(candidate.ToJSON())
	if err != nil {
		return err
	}
	msg := conference_utils.WebsocketMessage{
		Event: EventCandidate,
		Data:  string(data),
	}
	return session.Writer.WriteJSON(msg)
}

// Обработка входящих медиатреков (чтение RTP пакетов и ретрансляция)
func (s *WSHandler) handleTrack(session *PeerSession, t *webrtc.TrackRemote) {
	trackLocal, err := s.conference_usecase.AddTrack(&conference_utils.PeerConnection{
		PC:   session.PC,
		Conn: session.Writer,
	}, t)
	if err != nil {
		s.logger.Errorf("add track error: %v", err)
		return
	}
	defer s.conference_usecase.RemoveTrack(trackLocal)

	bufferSize := 1500
	if t.Kind() == webrtc.RTPCodecTypeAudio {
		bufferSize = 500
	}
	buf := make([]byte, bufferSize)

	for {
		n, _, err := t.Read(buf)
		if err != nil {
			s.logger.Infof("track read error or closed: %v", err)
			return
		}

		var pkt rtp.Packet
		if err := pkt.Unmarshal(buf[:n]); err != nil {
			s.logger.Errorf("unmarshal RTP packet error: %v", err)
			continue
		}

		// Отключаем расширения для видео, если не аудио
		if t.Kind() != webrtc.RTPCodecTypeAudio {
			pkt.Extension = false
			pkt.Extensions = nil
		}

		if err := trackLocal.WriteRTP(&pkt); err != nil {
			s.logger.Errorf("write RTP error: %v", err)
			return
		}
	}
}

// Обработка сообщений от клиента через WebSocket
func (s *WSHandler) handleWebsocketMessage(session *PeerSession, msg conference_utils.WebsocketMessage) error {
	switch msg.Event {
	case EventCandidate:
		var candidate webrtc.ICECandidateInit
		if err := json.Unmarshal([]byte(msg.Data), &candidate); err != nil {
			return err
		}
		return session.PC.AddICECandidate(candidate)

	case EventAnswer:
		var answer webrtc.SessionDescription
		if err := json.Unmarshal([]byte(msg.Data), &answer); err != nil {
			return err
		}
		return session.PC.SetRemoteDescription(answer)

	default:
		// Игнорируем неизвестные события, можно добавить логирование
		return nil
	}
}
