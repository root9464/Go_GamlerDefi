package conference_usecase

import (
	"sync"

	"github.com/gofiber/contrib/socketio"
	"github.com/pion/webrtc/v4"
	conference_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/entity"
	audiorecorder "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/usecase/audio_recorder"
	logger "github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type IConferenceUsecase interface {
	Disconect(ep *socketio.EventPayload)
	GetOrCreateRoom(roomID string, requestID string, conn *conference_entity.Connection) *conference_entity.Room
	CreateConnection(roomID string, pc *webrtc.PeerConnection, kws *socketio.Websocket) *conference_entity.Connection
	SetubWebRTC(conn *conference_entity.Connection, r *conference_entity.Room, requestID string)
	SignalPeerConnections(requestID string, roomID string)
	StartKeyFrameDispatcher()
}

type ConferenceUsecase struct {
	logger        *logger.Logger
	rooms         map[string]*conference_entity.Room
	roomsLock     sync.RWMutex
	serverRunning uint32
	bufferPool    sync.Pool

	audioRecorder *audiorecorder.AudioRecorder
}

func NewConferenceUsecase(logger *logger.Logger) IConferenceUsecase {
	return &ConferenceUsecase{
		rooms:         make(map[string]*conference_entity.Room),
		logger:        logger,
		serverRunning: 1,
		roomsLock:     sync.RWMutex{},
		audioRecorder: audiorecorder.NewAudioRecorder("../tmp/audiorecorder"),
		bufferPool: sync.Pool{
			New: func() any {
				buf := make([]byte, 1500)
				return &buf
			},
		},
	}
}
