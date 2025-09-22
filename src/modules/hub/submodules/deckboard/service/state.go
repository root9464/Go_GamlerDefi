package deckboard_service

import (
	"context"

	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IRepository interface {
	CreateState(ctx context.Context, data *deckboard_models.GameState) (*bson.ObjectID, error)
	RemoveState(ctx context.Context, stateID bson.ObjectID) error
	GetState(ctx context.Context, stateID bson.ObjectID) (*deckboard_models.GameState, error)
}

type StateManager struct {
	logger     *logger.Logger
	repository IRepository
}

type IStateManager interface {
	CreateState(ctx context.Context, data *deckboard_models.GameState) (*bson.ObjectID, error)
	RemoveState(ctx context.Context, stateID bson.ObjectID) error
	GetState(ctx context.Context, stateID bson.ObjectID) (*deckboard_models.GameState, error)
}

func NewStateManager(
	logger *logger.Logger,
	repository IRepository,
) IStateManager {
	return &StateManager{
		logger:     logger,
		repository: repository,
	}
}

func (sm *StateManager) CreateState(ctx context.Context, data *deckboard_models.GameState) (*bson.ObjectID, error) {
	stateID, err := sm.repository.CreateState(ctx, data)
	if err != nil {
		return nil, err
	}

	return stateID, nil
}

func (sm *StateManager) RemoveState(ctx context.Context, stateID bson.ObjectID) error {
	if err := sm.repository.RemoveState(ctx, stateID); err != nil {
		return err
	}

	return nil
}

func (sm *StateManager) GetState(ctx context.Context, stateID bson.ObjectID) (*deckboard_models.GameState, error) {
	state, err := sm.repository.GetState(ctx, stateID)
	if err != nil {
		return nil, err
	}

	return state, nil
}
