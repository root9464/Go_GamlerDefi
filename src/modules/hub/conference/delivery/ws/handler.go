package conference_ws

import (
	"github.com/gofiber/contrib/socketio"
	"github.com/pion/webrtc/v4"
	conference_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/entity"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/slog_logger"
)

type IConferenceUsecase interface {
	Disconect(ep *socketio.EventPayload)
	GetOrCreateRoom(roomID string, requestID string, conn *conference_entity.Connection) *conference_entity.Room
	CreateConnection(roomID string, pc *webrtc.PeerConnection, kws *socketio.Websocket) *conference_entity.Connection
	SetubWebRTC(conn *conference_entity.Connection, r *conference_entity.Room, requestID string)
	SignalPeerConnections(requestID string, roomID string)
}
type WSHandler struct {
	logger             *logger.Logger
	conference_usecase IConferenceUsecase
}

func NewWSHanler(
	logger *logger.Logger,
	conferenceUsecase IConferenceUsecase,
) *WSHandler {
	return &WSHandler{
		logger:             logger,
		conference_usecase: conferenceUsecase,
	}
}
