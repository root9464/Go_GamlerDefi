package game_session_hub

import (
	"context"
	"fmt"
	"sync"

	game_session_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/entity"
	game_session_registry "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/registry"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type Repository interface {
	Create(ctx context.Context, gameSession *game_session_entity.GameSession) error
	GetByID(ctx context.Context, id string) (*game_session_entity.GameSession, error)
	DeleteByID(ctx context.Context, id string) error
}

type Hub struct {
	logger     *logger.Logger
	repository Repository

	hub   map[string]*GameSession
	hubMU sync.RWMutex
}

func NewHub(logger *logger.Logger) *Hub {
	return &Hub{
		logger: logger,
		hub:    make(map[string]*GameSession),
	}
}

func (h *Hub) GetOrCreateSession(sessionID string, gameName string) (*GameSession, error) {
	h.hubMU.Lock()
	defer h.hubMU.Unlock()

	session, ok := h.hub[sessionID]
	if ok {
		if session.GameName != gameName {
			return nil, fmt.Errorf("сессия '%s' уже существует, но для другой игры ('%s')", sessionID, session.GameName)
		}
		return session, nil
	}

	gameLogic, err := game_session_registry.NewGame(gameName)
	if err != nil {
		return nil, err // Ошибка, если игра не найдена в реестре
	}

	session = &GameSession{
		ID:       sessionID,
		GameName: gameName, // <-- СОХРАНЯЕМ ИМЯ ИГРЫ
		Game:     gameLogic,
		Players:  make(map[string]*game_session_entity.Connection),
	}

	session.Game.Initialize(sessionID, session.SendToAll, session.SendToPlayer)
	h.hub[sessionID] = session
	h.logger.Infof("game session created: %s for game: %s", sessionID, gameName)
	return session, nil
}

func (h *Hub) FindSession(sessionID string) *GameSession {
	h.hubMU.RLock()
	defer h.hubMU.RUnlock()
	return h.hub[sessionID]
}
