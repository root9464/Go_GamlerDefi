package game_session_ws

import (
	"sync"

	"github.com/root9464/Go_GamlerDefi/src/config"
	conference_ws "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/handler"
	game_session_hub "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/hub.go"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type GameSessionHandler struct {
	hubManager    *game_session_hub.Hub
	logger        *logger.Logger
	uuidToSession map[string]*game_session_hub.GameSession
	mapMu         sync.RWMutex

	conferenceHandler *conference_ws.WSHandler
	config            *config.Config
}



func NewGameSessionHandler(hubManager *game_session_hub.Hub, logger *logger.Logger, conferenceHandler *conference_ws.WSHandler) *GameSessionHandler {
	return &GameSessionHandler{
		hubManager:        hubManager,
		logger:            logger,
		uuidToSession:     make(map[string]*game_session_hub.GameSession),
		conferenceHandler: conferenceHandler,
		config:            config,
	}
}
