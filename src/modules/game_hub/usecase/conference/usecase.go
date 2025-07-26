package conference_usecase

import (
	"sync"

	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type Room struct {
	ID          string
	peers       []*conference_utils.PeerConnection
	trackLocals map[string]*webrtc.TrackLocalStaticRTP
}

type IConferenceUsecase interface {
	AddTrack(peer *conference_utils.PeerConnection, t *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error)
	RemoveTrack(track *webrtc.TrackLocalStaticRTP)
	UpdatePeerTracks(peer *conference_utils.PeerConnection) error

	AddPeer(pc *webrtc.PeerConnection, conn *conference_utils.ThreadSafeWriter)
	SignalPeers() error
	DispatchKeyFrames()
}

type ConferenceUsecase struct {
	logger      *logger.Logger
	mu          sync.RWMutex
	trackLocals map[string]*webrtc.TrackLocalStaticRTP
	peers       []*conference_utils.PeerConnection
	rooms       map[string]*Room
}

func NewConferenceUsecase(logger *logger.Logger) IConferenceUsecase {
	return &ConferenceUsecase{
		logger:      logger,
		trackLocals: make(map[string]*webrtc.TrackLocalStaticRTP),
		mu:          sync.RWMutex{},
		rooms:       make(map[string]*Room),
	}
}
