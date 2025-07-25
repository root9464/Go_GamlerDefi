package conference_usecase

import (
	"sync"

	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type IConferenceUsecase interface {
	AddTrack(peer *conference_utils.PeerConnection, t *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error)
	RemoveTrack(track *webrtc.TrackLocalStaticRTP)
	UpdatePeerTracks(peer *conference_utils.PeerConnection) error

	AddPeer(pc *webrtc.PeerConnection, conn *conference_utils.ThreadSafeWriter)
	SignalPeers() error
	DispatchKeyFrames()
}

type ConferenceUsecase struct {
	logger              *logger.Logger
	mu                  sync.RWMutex
	trackLocals         map[string]*webrtc.TrackLocalStaticRTP
	trackOwners         map[string]*conference_utils.PeerConnection
	peers               []*conference_utils.PeerConnection
	signalingInProgress bool
}

func NewConferenceUsecase(logger *logger.Logger) IConferenceUsecase {
	return &ConferenceUsecase{
		logger:      logger,
		trackLocals: make(map[string]*webrtc.TrackLocalStaticRTP),
		trackOwners: make(map[string]*conference_utils.PeerConnection),
		mu:          sync.RWMutex{},
	}
}
