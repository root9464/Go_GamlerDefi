package game_config_repository

import (
	"context"
	"errors"

	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	CollectionName = "game_configs"
)

type GameConfigRepository struct {
	logger *logger.Logger
	db     *mongo.Database
}

func NewGameConfigRepository(logger *logger.Logger, db *mongo.Database) *GameConfigRepository {
	return &GameConfigRepository{logger: logger, db: db}
}

// func (r *GameSessionRepository) getCollection() *mongo.Collection {
// 	return r.db.Collection(CollectionName)
// }

func (r *GameConfigRepository) getCollection() *mongo.Collection {
	return r.db.Collection(CollectionName)
}

func (r *GameConfigRepository) GetByName(ctx context.Context, name string) (*GameConfigModel, error) {
	collection := r.getCollection()
	var config GameConfigModel

	filter := bson.M{"_id": name}
	if err := collection.FindOne(ctx, filter).Decode(&config); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warnf("game session with id '%s' not found", name)
			return nil, nil
		}
		r.logger.Errorf("failed to find game session with id '%s': %v", name, err)
		return nil, err
	}

	return &config, nil
}
