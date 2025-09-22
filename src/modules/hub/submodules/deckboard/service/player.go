package deckboard_service

import (
	"context"
	"fmt"

	deckboard_dto "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/dto"
	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type iPlayerRepository interface {
	UpdatePlayerInState(ctx context.Context, data *deckboard_models.Player, stateID bson.ObjectID) error
	GetAllPlayersFromState(ctx context.Context, stateID bson.ObjectID) ([]deckboard_models.Player, error)
	GetPlayerFromState(ctx context.Context, playerID string, stateID bson.ObjectID) (*deckboard_models.Player, error)
	RemovePlayerFromState(ctx context.Context, playerID string, stateID bson.ObjectID) error
	AddPlayerToState(ctx context.Context, data *deckboard_models.Player, stateID bson.ObjectID) error
}

type IPlayerManager interface {
	AddPlayer(ctx context.Context, data *deckboard_dto.AddPlayer) (*deckboard_models.Player, error)
	LeavePlayer(ctx context.Context, playerID string) error
	RemovePlayer(ctx context.Context, playerID string) error
	UpdatePlayerPosition(ctx context.Context, id string, position *deckboard_models.PlayerPosition) error
}

type PlayerManager struct {
	logger  *logger.Logger
	stateID bson.ObjectID

	repository iPlayerRepository
}

func NewPlayerManager(logger *logger.Logger, repository iPlayerRepository, stateID bson.ObjectID) IPlayerManager {
	return &PlayerManager{
		stateID:    stateID,
		logger:     logger,
		repository: repository,
	}
}

func (pm *PlayerManager) AddPlayer(ctx context.Context, data *deckboard_dto.AddPlayer) (*deckboard_models.Player, error) {
	existPlayer, err := pm.repository.GetPlayerFromState(ctx, data.UserID, pm.stateID)
	if err != nil {
		return nil, err
	}
	if existPlayer != nil {
		existPlayer.IsActive = true

		if err := pm.repository.UpdatePlayerInState(ctx, existPlayer, pm.stateID); err != nil {
			return nil, err
		}

		return existPlayer, nil
	}

	playerModel := deckboard_models.Player{
		ID:       data.UserID,
		Hand:     make([]deckboard_models.Card, 0),
		Position: deckboard_models.PlayerPosition{X: decimal.NewFromFloat32(0), Y: decimal.NewFromFloat32(0)},
		Metadata: data.Metadata,
		IsActive: true,
	}

	if err := pm.repository.AddPlayerToState(ctx, &playerModel, pm.stateID); err != nil {
		return nil, err
	}

	return &playerModel, nil
}

func (pm *PlayerManager) RemovePlayer(ctx context.Context, playerID string) error {
	if err := pm.repository.RemovePlayerFromState(ctx, playerID, pm.stateID); err != nil {
		return err
	}

	return nil
}

func (pm *PlayerManager) LeavePlayer(ctx context.Context, playerID string) error {
	existPlayer, err := pm.repository.GetPlayerFromState(ctx, playerID, pm.stateID)
	if err != nil {
		return err
	}

	if existPlayer == nil {
		return nil
	}

	existPlayer.IsActive = false
	if err := pm.repository.UpdatePlayerInState(ctx, existPlayer, pm.stateID); err != nil {
		return err
	}

	return nil
}

func (pm *PlayerManager) UpdatePlayerPosition(ctx context.Context, id string, position *deckboard_models.PlayerPosition) error {
	player, err := pm.repository.GetPlayerFromState(ctx, id, pm.stateID)
	if err != nil {
		return err
	}

	if player == nil {
		return fmt.Errorf("player not found")
	}

	player.Position = *position
	if err := pm.repository.UpdatePlayerInState(ctx, player, pm.stateID); err != nil {
		return err
	}

	return nil
}
