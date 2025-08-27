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

func (r *GameSessionRepository) GeAll(ctx context.Context) ([]*game_session_entity.GameSession, error) {
	collection := r.getCollection()
	var gameSessions []game_session_models.GameSession
	filter := bson.M{"status": bson.M{"$ne": game_session_entity.StatusFinished}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		r.logger.Errorf("failed to find game sessions: %v", err)
		return nil, err
	}
	err = cursor.All(ctx, &gameSessions)
	if err != nil {
		r.logger.Errorf("failed to decode game sessions: %v", err)
		return nil, err
	}
	var result []*game_session_entity.GameSession

	for _, gameSession := range gameSessions {
		result = append(result, r.Converter.ModelToEntity(&gameSession))
	}
	return result, nil
}

func (r *GameSessionRepository) CreateScheduled(ctx context.Context, gameSession *game_session_entity.GameSession, playerID, playerName, gameName string) error {
	collection := r.getCollection()
	gameSession.Status = game_session_entity.StatusScheduled
	gameSession.GameName = gameName
	gameSession.Participants = []game_session_entity.Player{{
		ID:   playerID,
		Name: playerName,
	}}
	_, err := collection.InsertOne(ctx, r.Converter.EntityToModel(gameSession))
	if err != nil {
		r.logger.Errorf("failed to insert game session: %v", err)
		return err
	}
	return nil
}

func (r *GameSessionRepository) GetScheduledByID(ctx context.Context, id string) (*game_session_entity.GameSession, error) {
	collection := r.getCollection()
	hexID, err := r.Converter.StringToID(id)
	if err != nil {
		return nil, err
	}

	var gameSession game_session_models.GameSession
	filter := bson.M{
		"_id":    hexID,
		"status": bson.M{"$ne": game_session_entity.StatusFinished},
	}

	err = collection.FindOne(ctx, filter).Decode(&gameSession)
	if err != nil {
		r.logger.Errorf("failed to find game session: %v", err)
		return nil, err
	}

	return r.Converter.ModelToEntity(&gameSession), nil
}

func (r *GameSessionRepository) AddParticipant(ctx context.Context, id string, participantID string) error {
	collection := r.getCollection()
	update := bson.M{
		"$addToSet": bson.M{"participants": participantID},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		r.logger.Errorf("failed to update game session: %v", err)
		return err
	}

	return nil
}

func (r *GameSessionRepository) ActivateStatus(ctx context.Context, id string) error {
	collection := r.getCollection()
	update := bson.M{
		"$set": bson.M{"status": game_session_entity.StatusActive},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		r.logger.Errorf("failed to update game session: %v", err)
		return err
	}

	return nil
}
