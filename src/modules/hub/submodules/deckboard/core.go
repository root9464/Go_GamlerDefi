package deckboard_core

import (
	"context"

	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
	deckboard_dto "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/dto"
	game_state_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/infrastucture/repository"
	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
	deckboard_service "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/service"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GameConfig struct {
	HostID    string
	Decks     []deckboard_models.Deck
	SessionID string

	DiceCount int `json:"dice_count"`
	DiceFaces int `json:"dice_faces"`
}

type DataProvider interface {
	GetCoreGameConfig() *GameConfig // Метод для получения core-конфигурации игры
	GetSessionID() string           // Метод для получения ID текущей сессии
	GetHostID() string              // Метод для получения ID хоста текущей сессии
}

type DeckboardTemplate struct {
	config *GameConfig
	logger *logger.Logger
	db     *mongo.Database

	SendToAll            func(message any)
	SendToPlayer         func(userID string, message any) error
	BroadcastToAllExcept func(excludeUserID string, message any)

	diceManager   deckboard_service.IDiceManager
	deckManager   deckboard_service.IDeckManager
	playerManager deckboard_service.IPlayerManager
	stateManager  deckboard_service.IStateManager

	stateRepository game_state_repository.IGameStateRepository
}

// NewDeckboardTemplate - конструктор для создания нового экземпляра DeckboardTemplate.
// устанавливает, не зависящие от конфигуратора, зависимости
func NewDeckboardTemplate(
	db *mongo.Database,
	logger *logger.Logger,
	sendToAll func(message any),
	sendToPlayer func(userID string, message any) error,
	broadcastToAllExcept func(excludeUserID string, message any),
) *DeckboardTemplate {
	t := &DeckboardTemplate{
		logger:               logger,
		db:                   db,
		SendToAll:            sendToAll,
		SendToPlayer:         sendToPlayer,
		BroadcastToAllExcept: broadcastToAllExcept,
	}

	t.stateRepository = game_state_repository.NewRepository(t.db, t.logger)
	t.diceManager = deckboard_service.NewDiceManager(t.logger)

	return t
}

func (t *DeckboardTemplate) AddPlayer(ctx context.Context, userID string, metadata map[string]any) {
	data := deckboard_dto.AddPlayer{
		UserID:   userID,
		Metadata: metadata,
	}
	newPlayer, err := t.playerManager.AddPlayer(ctx, &data)
	if err != nil {
		t.SendFullStateToPlayer(userID)
		return
	}

	t.SendFullStateToPlayer(userID)
	event := game_session_contracts.Action{
		Type:    "player_joined",
		Payload: deckboard_models.PlayerJoinedPayload{Player: newPlayer},
	}
	t.SendToAll(event)
}

func (t *Template) RemovePlayer(userID string) {
	if t.PlayerManager.RemovePlayer(userID) {
		event := game_session_contracts.Action{
			Type:    "player_left",
			Payload: deckboard_models.PlayerLeftPayload{PlayerID: userID},
		}
		t.SendToAll(event)
	}
}

// Initialize настраивает специфическую для игры конфигурацию и менеджеры
func (t *DeckboardTemplate) Inizialize(
	ctx context.Context,
	provider DataProvider,
) {
	t.logger.Infof("Initializing the DeckboardTemplate for a session: %s", provider.GetSessionID())

	t.config = provider.GetCoreGameConfig()
	t.config.SessionID = provider.GetSessionID()
	t.config.HostID = provider.GetHostID()

	t.stateManager = deckboard_service.NewStateManager(t.logger, t.stateRepository)

	data := deckboard_models.NewGameState(t.config.Decks, t.config.DiceCount, t.config.DiceFaces, t.config.HostID)
	stateID, err := t.stateManager.CreateState(ctx, data)
	if err != nil {
		t.logger.Errorf("failed to create base game state for a session %s: %v", t.config.SessionID, err)
	}

	t.deckManager = deckboard_service.NewDeckManager(t.logger, t.stateRepository, *stateID)
	t.playerManager = deckboard_service.NewPlayerManager(t.logger, t.stateRepository, *stateID)
}
