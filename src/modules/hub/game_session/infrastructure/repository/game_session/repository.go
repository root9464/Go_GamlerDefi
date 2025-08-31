// src/modules/hub/game_session/infrastructure/repository/game_session/repository.go
package game_session_repository

import (
	"context"
	"errors"
	"time"

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

func (r *GameSessionRepository) Create(ctx context.Context, sess *game_session_entity.GameSession) (*game_session_entity.GameSession, error) {
	collection := r.getCollection()

	model := r.Converter.EntityToModel(sess)
	model.TimeToStart = time.Now()

	_, err := collection.InsertOne(ctx, model)
	if err != nil {
		r.logger.Errorf("failed to create game session: %v", err)
		return nil, err
	}

	r.logger.Infof("Created new game session with ID: %s", sess.ID)

	return r.Converter.ModelToEntity(model), nil
}

func (r *GameSessionRepository) GetByID(ctx context.Context, id string) (*game_session_entity.GameSession, error) {
	collection := r.getCollection()

	var model game_session_models.GameSession
	filter := bson.M{"_id": id}

	err := collection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warnf("game session with id '%s' not found", id)
			return nil, nil
		}
		r.logger.Errorf("failed to find game session with id '%s': %v", id, err)
		return nil, err
	}

	return r.Converter.ModelToEntity(&model), nil
}

func (r *GameSessionRepository) GetAll(ctx context.Context) ([]*game_session_entity.GameSession, error) {
	collection := r.getCollection()
	filter := bson.M{}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		r.logger.Errorf("failed to find game sessions: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var models []*game_session_models.GameSession
	if err = cursor.All(ctx, &models); err != nil {
		r.logger.Errorf("failed to decode game sessions: %v", err)
		return nil, err
	}

	entities := make([]*game_session_entity.GameSession, 0, len(models))
	for _, model := range models {
		entities = append(entities, r.Converter.ModelToEntity(model))
	}

	return entities, nil
}
