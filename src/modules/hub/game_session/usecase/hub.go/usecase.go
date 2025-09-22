package game_session_hub

import (
	"context"
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
	h.logger.Infof("starting session activation: sessionID=%s, userID=%s, gameName=%s", sessionID, userID, gameName)

	h.hubMU.Lock()
	defer func() {
		h.hubMU.Unlock()
		h.logger.Infof("mutex unlocked for session: sessionID=%s", sessionID)
	}()

	h.logger.Infof("parsing sessionID to uint: sessionID=%s", sessionID)
	uintID, err := strconv.ParseUint(sessionID, 10, 32)
	if err != nil {
		h.logger.Errorf("failed to parse sessionID: sessionID=%s, error=%v", sessionID, err)
		return nil, err
	}
	h.logger.Infof("sessionID parsed successfully: sessionID=%s, uintID=%d", sessionID, uintID)

	h.logger.Infof("fetching session from trash repository: sessionID=%d", uintID)
	session, err := h.trashRepo.GetSessionByID(ctx, uint(uintID))
	if err != nil {
		h.logger.Errorf("failed to get session from trash repo: sessionID=%d, error=%v", uintID, err)
		return nil, err
	}
	h.logger.Infof("session retrieved from trash repo: sessionID=%d", uintID)

	// if session.HostID != userID {
	// 	ok, err := h.trashRepo.IsUserHasAccess(ctx, uint(uintID), userID)
	// 	if err != nil || !ok {
	// 		return nil, errors.New("access denied")
	// 	}
	// }

	if activeSession, ok := h.hub[sessionID]; ok {
		h.logger.Infof("session already active in hub, returning existing session: sessionID=%s", sessionID)
		return activeSession, nil
	}
	h.logger.Infof("session not found in hub, creating new one: sessionID=%s", sessionID)

	h.logger.Infof("checking if game session exists in repository: sessionID=%s", sessionID)
	gameSession, err := h.repository.GetByID(ctx, sessionID)
	if err != nil {
		h.logger.Errorf("failed to get game session from repository: sessionID=%s, error=%v", sessionID, err)
		return nil, err
	}

	if gameSession == nil {
		h.logger.Infof("game session not found, creating new session: sessionID=%s", sessionID)

		h.logger.Infof("fetching session host from trash repo: sessionID=%d", uintID)
		host, err := h.trashRepo.GetSessionHost(ctx, uint(uintID))
		if err != nil {
			h.logger.Errorf("failed to get session host: sessionID=%d, error=%v", uintID, err)
			return nil, err
		}
		h.logger.Infof("session host retrieved: sessionID=%d, host=%s", uintID, host)

		gameSession = &game_session_entity.GameSession{
			ID:       sessionID,
			GameName: gameName,
			HostID:   host,
		}

		h.logger.Infof("creating new game session in repository: sessionID=%s, gameName=%s, host=%s", sessionID, gameName, host)
		_, err = h.repository.Create(ctx, gameSession)
		if err != nil {
			h.logger.Errorf("failed to create game session: sessionID=%s, error=%v", sessionID, err)
			return nil, err
		}
		h.logger.Infof("game session created successfully: sessionID=%s", sessionID)
	} else {
		h.logger.Infof("existing game session found: sessionID=%s, gameName=%s", sessionID, gameSession.GameName)
	}

	h.logger.Infof("initializing game logic for game: %s", gameSession.GameName)
	gameLogic, err := game_session_registry.NewGame(gameSession.GameName, *h.gameConfigRepository)
	if err != nil {
		h.logger.Errorf("failed to initialize game logic: gameName=%s, error=%v", gameSession.GameName, err)
		return nil, err
	}
	h.logger.Infof("game logic initialized successfully: gameName=%s", gameSession.GameName)

	activeSession := &GameSession{
		ID:       sessionID,
		GameName: gameSession.GameName,
		Game:     gameLogic,
		Players:  make(map[string]*game_session_entity.Connection),
		HostID:   session.HostID,
	}

	h.logger.Infof("initializing game with host: sessionID=%s, hostID=%s", sessionID, gameSession.HostID)
	activeSession.Game.Initialize(gameSession.HostID, activeSession.SendToAll, activeSession.SendToPlayer, activeSession.BroadcastToAllExcept)

	h.hub[sessionID] = activeSession
	h.logger.Infof("session added to hub: sessionID=%s, gameName=%s", sessionID, gameSession.GameName)

	h.logger.Infof("session activation completed successfully: sessionID=%s", sessionID)
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
