package game_session_hub

import (
	"context"
	"errors"
	"strconv"
	"sync"

	game_session_dto "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/dto"
	game_session_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/entity"
	game_session_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/infrastructure/repository/game_session"
	trash_repo "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/infrastructure/trash/repository"
	game_session_registry "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/registry"
	game_config_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/repository"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type Hub struct {
	logger               *logger.Logger
	repository           *game_session_repository.GameSessionRepository
	trashRepo            *trash_repo.TrashRepository
	gameConfigRepository *game_config_repository.GameConfigRepository

	hub   map[string]*GameSession
	hubMU sync.RWMutex
}

func NewHub(
	logger *logger.Logger,
	repository *game_session_repository.GameSessionRepository,
	trashRepo *trash_repo.TrashRepository,
	gameConfigRepository *game_config_repository.GameConfigRepository,
) *Hub {
	return &Hub{
		logger:               logger,
		repository:           repository,
		hub:                  make(map[string]*GameSession),
		trashRepo:            trashRepo,
		gameConfigRepository: gameConfigRepository,
	}
}

func (h *Hub) FindSession(sessionID string) *GameSession {
	h.hubMU.RLock()
	defer h.hubMU.RUnlock()
	return h.hub[sessionID]
}

func (h *Hub) ActiveteSession(ctx context.Context, sessionID, userID, gameName string) (*GameSession, error) {
	h.hubMU.Lock()
	defer h.hubMU.Unlock()

	uintID, err := strconv.ParseUint(sessionID, 10, 32)
	if err != nil {
		return nil, err
	}

	session, err := h.trashRepo.GetSessionByID(ctx, uint(uintID))
	if err != nil {
		return nil, err
	}

	if session.HostID != userID {
		ok, err := h.trashRepo.IsUserHasAccess(ctx, uint(uintID), userID)
		if err != nil || !ok {
			return nil, errors.New("access denied")
		}
	}

	if activeSession, ok := h.hub[sessionID]; ok {
		return activeSession, nil
	}

	gameSession, err := h.repository.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if gameSession == nil {
		host, err := h.trashRepo.GetSessionHost(ctx, uint(uintID))
		if err != nil {
			return nil, err
		}

		gameSession = &game_session_entity.GameSession{
			ID:       sessionID,
			GameName: gameName,
			HostID:   host,
		}

		_, err = h.repository.Create(ctx, gameSession)
		if err != nil {
			return nil, err
		}
	}

	gameLogic, err := game_session_registry.NewGame(gameSession.GameName, *h.gameConfigRepository)
	if err != nil {
		return nil, err
	}

	activeSession := &GameSession{
		ID:       sessionID,
		GameName: gameSession.GameName,
		Game:     gameLogic,
		Players:  make(map[string]*game_session_entity.Connection),
	}

	activeSession.Game.Initialize(gameSession.HostID, activeSession.SendToAll, activeSession.SendToPlayer, activeSession.BroadcastToAllExcept)
	h.hub[sessionID] = activeSession

	h.logger.Infof("game session activated: %s for game: %s", sessionID, gameSession.GameName)
	return activeSession, nil
}

func (h *Hub) GetAllSessions(ctx context.Context) ([]*game_session_entity.GameSession, error) {
	sessions, err := h.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	sessionsInfo := make([]*game_session_dto.SessionInfo, len(sessions))
	for i, s := range sessions {
		sessionsInfo[i] = &game_session_dto.SessionInfo{
			ID:       s.ID,
			GameName: s.GameName,
			HostID:   s.HostID,
		}
	}

	return sessions, nil
}
