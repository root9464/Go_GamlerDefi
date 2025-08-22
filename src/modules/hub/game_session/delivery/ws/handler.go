package game_session_ws

import (
	"sync"

	game_session_hub "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/hub.go"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type GameSessionHandler struct {
	hubManager    *game_session_hub.Hub
	logger        *logger.Logger
	uuidToSession map[string]*game_session_hub.GameSession
	mapMu         sync.RWMutex
}

func NewGameSessionHandler(hubManager *game_session_hub.Hub, logger *logger.Logger) *GameSessionHandler {
	return &GameSessionHandler{
		hubManager:    hubManager,
		logger:        logger,
		uuidToSession: make(map[string]*game_session_hub.GameSession),
	}
}
