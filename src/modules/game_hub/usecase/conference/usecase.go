package conference_usecase

import (
	"sync"
	"time"

	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type Hub struct {
	ID          string
	peers       []*conference_utils.PeerConnection
	trackLocals map[string]*webrtc.TrackLocalStaticRTP
	mu          sync.RWMutex
}

type IConferenceUsecase interface {
	AddTrack(pc *conference_utils.PeerConnection, t *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error)
	RemoveTrack(t *webrtc.TrackLocalStaticRTP, pc *conference_utils.PeerConnection)
	UpdatePeerTracks(pc *conference_utils.PeerConnection) error

	AddPeer(wpc *webrtc.PeerConnection, pc *conference_utils.PeerConnection)
	SignalPeers(pc *conference_utils.PeerConnection) error
	DispatchKeyFrames(hubID string)

	JoinHub(pc *conference_utils.PeerConnection) error
	LeaveHub(pc *conference_utils.PeerConnection) error
}

type ConferenceUsecase struct {
	logger     *logger.Logger
	hubs       map[string]*Hub
	hubTickers map[string]*time.Ticker
	mu         sync.RWMutex
}

func NewConferenceUsecase(logger *logger.Logger) IConferenceUsecase {
	return &ConferenceUsecase{
		logger:     logger,
		hubs:       make(map[string]*Hub),
		hubTickers: make(map[string]*time.Ticker),
		mu:         sync.RWMutex{},
	}
}
