package game_session_hub

import (
	"context"
	"errors"
	"sync"

	game_session_dto "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/dto"
	game_session_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/entity"
	game_session_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/infrastructure/repository/game_session"
	game_session_registry "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/registry"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type Hub struct {
	logger     *logger.Logger
	repository *game_session_repository.GameSessionRepository

	hub   map[string]*GameSession
	hubMU sync.RWMutex
}

func NewHub(logger *logger.Logger, repository *game_session_repository.GameSessionRepository) *Hub {
	return &Hub{
		logger:     logger,
		repository: repository,
		hub:        make(map[string]*GameSession),
	}
}

func (h *Hub) FindSession(sessionID string) *GameSession {
	h.hubMU.RLock()
	defer h.hubMU.RUnlock()
	return h.hub[sessionID]
}

func (h *Hub) ActiveteSession(ctx context.Context, sessionID string) (*GameSession, error) {
	h.hubMU.Lock()
	defer h.hubMU.Unlock()

	if activeSession, ok := h.hub[sessionID]; ok {
		return activeSession, nil
	}

	gameSession, err := h.repository.GetScheduledByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	gameLogic, err := game_session_registry.NewGame(gameSession.GameName)
	if err != nil {
		return nil, err
	}

	activeSession := &GameSession{
		ID:       sessionID,
		GameName: gameSession.GameName,
		Game:     gameLogic,
		Players:  make(map[string]*game_session_entity.Connection),
	}

	activeSession.Game.Initialize(activeSession.SendToAll, activeSession.SendToPlayer, activeSession.BroadcastToAllExcept)
	h.hub[sessionID] = activeSession

	h.logger.Infof("game session activated: %s for game: %s", sessionID, gameSession.GameName)
	return activeSession, nil
}

func (h *Hub) CreateSession(ctx context.Context, playerID, playerName, gameName string) error {
	h.logger.Infof("create session for player: %s and game: %s", playerID, gameName)

	if !game_session_registry.IsGameRegistered(gameName) {
		return errors.New("game not registered")
	}
	gameSession := game_session_entity.GameSession{
		HostID: playerID,
	}

	return h.repository.CreateScheduled(ctx, &gameSession, playerID, playerName, gameName)
}

func (h *Hub) GetAllSessions(ctx context.Context) ([]*game_session_entity.GameSession, error) {
	sessions, err := h.repository.GeAll(ctx)
	if err != nil {
		return nil, err
	}

	sessionsInfo := make([]*game_session_dto.SessionInfo, len(sessions))
	for i, s := range sessions {
		sessionsInfo[i] = &game_session_dto.SessionInfo{
			ID:          s.ID,
			GameName:    s.GameName,
			PlayerCount: len(s.Participants),
			HostID:      s.HostID,
		}
	}

	return sessions, nil
}
