package game_session_repository

import (
	"context"

	game_session_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/entity"
	game_session_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/infrastructure/repository/models"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	CollectionName = "game_sessions"
)

type GameSessionRepository struct {
	Converter *Converter
	logger    *logger.Logger
	db        *mongo.Database
}

func NewGameSessionRepository(logger *logger.Logger, db *mongo.Database) *GameSessionRepository {
	return &GameSessionRepository{
		Converter: NewConverter(),
		logger:    logger,
		db:        db,
	}
}

func (r *GameSessionRepository) getCollection() *mongo.Collection {
	return r.db.Collection(CollectionName)
}

func (r *GameSessionRepository) Create(ctx context.Context, gameSession *game_session_entity.GameSession) error {
	collection := r.getCollection()
	_, err := collection.InsertOne(ctx, r.Converter.EntityToModel(gameSession))
	if err != nil {
		r.logger.Errorf("failed to insert game session: %v", err)
		return err
	}
	return err
}

func (r *GameSessionRepository) GetByID(ctx context.Context, id string) (*game_session_entity.GameSession, error) {
	collection := r.getCollection()
	var gameSession game_session_models.GameSession
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&gameSession)
	if err != nil {
		r.logger.Errorf("failed to find game session: %v", err)
		return nil, err
	}
	return r.Converter.ModelToEntity(&gameSession), nil
}

func (r *GameSessionRepository) DeleteByID(ctx context.Context, id string) error {
	collection := r.getCollection()
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		r.logger.Errorf("failed to delete game session: %v", err)
		return err
	}
	return nil
}
