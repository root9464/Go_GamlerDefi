package conference_ws

import (
	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type IConferenceUsecase interface {
	AddTrack(pc *conference_utils.PeerConnection, t *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error)
	RemoveTrack(t *webrtc.TrackLocalStaticRTP, pc *conference_utils.PeerConnection)

	AddPeer(pc *conference_utils.PeerConnection)
	SignalPeers(pc *conference_utils.PeerConnection) error

	JoinHub(hubID, userID string) error
	// LeaveHub(pc *conference_utils.PeerConnection) error
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
