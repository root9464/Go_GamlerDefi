package conference_ws

import (
	"sync"

	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type IConferenceUsecase interface {
	AddPeer(pc *webrtc.PeerConnection, conn *conference_utils.ThreadSafeWriter)
	SignalPeers() error

	AddTrack(peer *conference_utils.PeerConnection, t *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error)
	RemoveTrack(track *webrtc.TrackLocalStaticRTP)
}

type WSHandler struct {
	logger             *logger.Logger
	conference_usecase IConferenceUsecase
	peers              sync.Map
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
