package game_config_repository

import (
	"context"
	"errors"

	game_config_model "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/model"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (r *GameConfigRepository) getCollection() *mongo.Collection {
	return r.db.Collection(CollectionName)
}

func (r *GameConfigRepository) GetByName(ctx context.Context, name string) (*game_config_model.GameConfigModel, error) {
	collection := r.getCollection()
	var config game_config_model.GameConfigModel

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

func (r *GameConfigRepository) Upsert(ctx context.Context, config *game_config_model.GameConfigModel) error {
	collection := r.getCollection()

	if config.ID == "" {
		return errors.New("config ID (game name) cannot be empty")
	}
	filter := bson.M{"_id": config.ID}
	update := bson.M{"$set": config}
	opts := options.UpdateOne().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// ===================================================================
func (r *GameConfigRepository) FindByID(ctx context.Context, id string) (*game_config_model.GameConfigModel, error) {
	collection := r.getCollection()
	var config game_config_model.GameConfigModel

	filter := bson.M{"_id": id}
	if err := collection.FindOne(ctx, filter).Decode(&config); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warnf("game session with id '%s' not found", id)
			return nil, nil
		}
		r.logger.Errorf("failed to find game session with id '%s': %v", id, err)
		return nil, err
	}

	return &config, nil
}

func (r *GameConfigRepository) Create(ctx context.Context, config *game_config_model.GameConfigModel) error {
	collection := r.getCollection()
	_, err := collection.InsertOne(ctx, config)
	return err
}

func (r *GameConfigRepository) UpdateSettings(ctx context.Context, id string, settings primitive.M) error {
	collection := r.getCollection()
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"settings": settings}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document matched the filter")
	}
	return nil
}
