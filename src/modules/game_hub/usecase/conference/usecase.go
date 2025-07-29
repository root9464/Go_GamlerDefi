package conference_usecase

import (
	"sync"

	"github.com/gofiber/contrib/socketio"
	"github.com/pion/webrtc/v4"
	hub_entity "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/entity"
	logger "github.com/root9464/Go_GamlerDefi/src/packages/lib/slog_logger"
)

type IConferenceUsecase interface {
	Disconect(ep *socketio.EventPayload)
	GetOrCreateRoom(roomID string, requestID string, conn *hub_entity.Connection) *hub_entity.Room
	CreateConnection(roomID string, pc *webrtc.PeerConnection, kws *socketio.Websocket) *hub_entity.Connection
	SetubWebRTC(conn *hub_entity.Connection, r *hub_entity.Room, requestID string)
	SignalPeerConnections(requestID string, roomID string)
	StartKeyFrameDispatcher()
}

type ConferenceUsecase struct {
	logger        *logger.Logger
	rooms         map[string]*hub_entity.Room
	roomsLock     sync.RWMutex
	serverRunning uint32
	bufferPool    sync.Pool
}

func NewConferenceUsecase(logger *logger.Logger) IConferenceUsecase {
	return &ConferenceUsecase{
		rooms:         make(map[string]*hub_entity.Room),
		logger:        logger,
		serverRunning: 1,
		roomsLock:     sync.RWMutex{},
		bufferPool: sync.Pool{
			New: func() any {
				buf := make([]byte, 1500)
				return &buf
			},
		},
	}
}
