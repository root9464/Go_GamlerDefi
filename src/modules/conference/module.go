package conference_module

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pion/webrtc/v4"
	conference_ws_handler "github.com/root9464/Go_GamlerDefi/src/modules/conference/delivery/ws"
	conference_peer_service "github.com/root9464/Go_GamlerDefi/src/modules/conference/services/peer"
	conference_track_service "github.com/root9464/Go_GamlerDefi/src/modules/conference/services/track"
	"github.com/root9464/Go_GamlerDefi/src/modules/conference/utils"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type ConferenceModule struct {
	peerService  conference_peer_service.IPeerService
	trackService conference_track_service.ITrackService
	wsHandler    *conference_ws_handler.WSHandler

	logger      *logger.Logger
	trackLocals map[string]*webrtc.TrackLocalStaticRTP
	trackOwners map[string]*conference_utils.PeerConnection
}

func NewConferenceModule(logger *logger.Logger) *ConferenceModule {
	return &ConferenceModule{
		logger:      logger,
		trackLocals: make(map[string]*webrtc.TrackLocalStaticRTP),
		trackOwners: make(map[string]*conference_utils.PeerConnection),
	}
}

func (m *ConferenceModule) Init() {
	m.trackService = conference_track_service.NewTrackService(m.logger, m.trackLocals, m.trackOwners)
	m.peerService = conference_peer_service.NewPeerService(m.logger, m.trackLocals, m.trackOwners, m.trackService)
	m.wsHandler = conference_ws_handler.NewWSHanler(m.logger, m.peerService, m.trackService)

	go func() {
		for range time.Tick(3 * time.Second) {
			m.peerService.DispatchKeyFrames()
		}
	}()
}

func (m *ConferenceModule) InitRoutes(app fiber.Router) {
	m.Init()

	fmt.Println("ConMod: ", m)
	app.Get("/websocket", m.wsHandler.HandleWebSocketFiber)
}
