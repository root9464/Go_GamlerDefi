package game_state_repository

import (
	"context"
	"fmt"

	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type iUserRepository interface {
	UpdatePlayerInState(ctx context.Context, data *deckboard_models.Player, stateID bson.ObjectID) error
	GetAllPlayersFromState(ctx context.Context, stateID bson.ObjectID) ([]deckboard_models.Player, error)
	GetPlayerFromState(ctx context.Context, playerID string, stateID bson.ObjectID) (*deckboard_models.Player, error)
	RemovePlayerFromState(ctx context.Context, playerID string, stateID bson.ObjectID) error
	AddPlayerToState(ctx context.Context, data *deckboard_models.Player, stateID bson.ObjectID) error
}

func (r *GameStateRepository) AddPlayerToState(ctx context.Context, data *deckboard_models.Player, stateID bson.ObjectID) error {
	filter := bson.M{"_id": stateID}
	update := bson.M{
		"$push": bson.M{
			"players": data,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to add player to game state: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("game state with ID %s not found", stateID.Hex())
	}

	r.logger.Infof("added player %s to game state %s", data.ID, stateID.Hex())
	return nil
}

func (r *GameStateRepository) RemovePlayerFromState(ctx context.Context, playerID string, stateID bson.ObjectID) error {
	filter := bson.M{"_id": stateID}
	update := bson.M{
		"$pull": bson.M{
			"players": bson.M{"id": playerID},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to remove player from game state: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("game state with ID %s not found", stateID.Hex())
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("player %s not found in game state %s", playerID, stateID.Hex())
	}

	r.logger.Infof("removed player %s from game state %s", playerID, stateID.Hex())
	return nil
}

func (r *GameStateRepository) GetPlayerFromState(ctx context.Context, playerID string, stateID bson.ObjectID) (*deckboard_models.Player, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": stateID},
		},
		{
			"$unwind": "$players",
		},
		{
			"$match": bson.M{"players.id": playerID},
		},
		{
			"$replaceRoot": bson.M{"newRoot": "$players"},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("failed to find player: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var player deckboard_models.Player
		if err := cursor.Decode(&player); err != nil {
			return nil, err
		}
		return &player, nil
	}

	return nil, nil
}

func (r *GameStateRepository) GetAllPlayersFromState(ctx context.Context, stateID bson.ObjectID) ([]deckboard_models.Player, error) {
	opts := options.FindOne().SetProjection(bson.M{"players": 1})

	var result struct {
		Players []deckboard_models.Player `bson:"players"`
	}

	filter := bson.M{"_id": stateID}
	err := r.collection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("game state with ID %s not found", stateID.Hex())
		}
		r.logger.Errorf("failed to find players: %v", err)
		return nil, err
	}

	return result.Players, nil
}

func (r *GameStateRepository) UpdatePlayerInState(ctx context.Context, data *deckboard_models.Player, stateID bson.ObjectID) error {
	filter := bson.M{
		"_id":        stateID,
		"players.id": data.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"players.$.position":  data.Position,
			"players.$.hand":      data.Hand,
			"players.$.metadata":  data.Metadata,
			"players.$.is_active": data.IsActive,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to update player: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("player %s not found in game state %s", data.ID, stateID.Hex())
	}

	r.logger.Infof("updated player %s in game state %s", data.ID, stateID.Hex())
	return nil
}
