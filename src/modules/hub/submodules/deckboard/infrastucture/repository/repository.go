package game_state_repository

import (
	"context"
	"fmt"

	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	CollectionName = "games_state"
)

type IGameStateRepository interface {
	iDeckRepository
	iUserRepository

	CreateState(ctx context.Context, data *deckboard_models.GameState) (*bson.ObjectID, error)
	RemoveState(ctx context.Context, stateID bson.ObjectID) error
	GetState(ctx context.Context, stateID bson.ObjectID) (*deckboard_models.GameState, error)
}

type GameStateRepository struct {
	db     *mongo.Database
	logger *logger.Logger

	collection *mongo.Collection
}

func NewRepository(db *mongo.Database, logger *logger.Logger) IGameStateRepository {
	return &GameStateRepository{
		db:         db,
		logger:     logger,
		collection: db.Collection(CollectionName),
	}
}

func (r *GameStateRepository) CreateState(ctx context.Context, data *deckboard_models.GameState) (*bson.ObjectID, error) {
	if _, err := r.collection.InsertOne(ctx, data); err != nil {
		r.logger.Errorf("failed to create game state: %v", err)
		return nil, err
	}

	r.logger.Infof("created new game state with ID: %s", data.ID)
	return &data.ID, nil
}

func (r *GameStateRepository) RemoveState(ctx context.Context, stateID bson.ObjectID) error {
	filter := bson.M{"_id": stateID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Errorf("failed to remove game state: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("game state with ID %s not found", stateID.Hex())
	}

	r.logger.Infof("successfully removed game state with ID: %s", stateID.Hex())
	return nil
}

func (r *GameStateRepository) GetState(ctx context.Context, stateID bson.ObjectID) (*deckboard_models.GameState, error) {
	filter := bson.M{"_id": stateID}

	var gameState deckboard_models.GameState
	err := r.collection.FindOne(ctx, filter).Decode(&gameState)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Errorf("failed to get game state: %v", err)
		return nil, err
	}

	return &gameState, nil
}
