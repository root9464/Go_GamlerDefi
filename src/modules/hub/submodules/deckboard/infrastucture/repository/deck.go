package game_state_repository

import (
	"context"
	"fmt"

	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type iDeckRepository interface {
	AddDecksToState(ctx context.Context, decks []*deckboard_models.Deck, stateID bson.ObjectID) error
	GetDeckFromState(ctx context.Context, stateID bson.ObjectID, deckID bson.ObjectID) (*deckboard_models.Deck, error)
	GetAllDecksFromState(ctx context.Context, stateID bson.ObjectID) ([]deckboard_models.Deck, error)
	RemoveDeckFromState(ctx context.Context, stateID bson.ObjectID, deckID bson.ObjectID) error
	UpdateDeckInState(ctx context.Context, stateID bson.ObjectID, deckID bson.ObjectID, updateData *deckboard_models.Deck) error

	RemoveCardFromDeck(ctx context.Context, stateID, deckID, cardID bson.ObjectID) error
}

func (r *GameStateRepository) AddDecksToState(ctx context.Context, decks []*deckboard_models.Deck, stateID bson.ObjectID) error {
	filter := bson.M{"_id": stateID}
	update := bson.M{
		"$push": bson.M{
			"decks": bson.M{
				"$each": decks,
			},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to add decks to game state: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("game state with ID %s not found", stateID.Hex())
	}

	r.logger.Infof("added %d decks to game state %s", len(decks), stateID.Hex())
	return nil
}

func (r *GameStateRepository) GetAllDecksFromState(ctx context.Context, stateID bson.ObjectID) ([]deckboard_models.Deck, error) {
	opts := options.FindOne().SetProjection(bson.M{"decks": 1})

	var result struct {
		Decks []deckboard_models.Deck `bson:"decks"`
	}

	filter := bson.M{"_id": stateID}
	err := r.collection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("game state with ID %s not found", stateID.Hex())
		}
		r.logger.Errorf("failed to find decks: %v", err)
		return nil, err
	}

	return result.Decks, nil
}

func (r *GameStateRepository) GetDeckFromState(ctx context.Context, stateID bson.ObjectID, deckID bson.ObjectID) (*deckboard_models.Deck, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": stateID},
		},
		{
			"$unwind": "$decks",
		},
		{
			"$match": bson.M{"decks.id": deckID},
		},
		{
			"$replaceRoot": bson.M{"newRoot": "$decks"},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("failed to find deck: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var deck deckboard_models.Deck
		if err := cursor.Decode(&deck); err != nil {
			return nil, err
		}
		return &deck, nil
	}

	return nil, fmt.Errorf("deck %s not found in game state %s", deckID, stateID.Hex())
}

func (r *GameStateRepository) RemoveDeckFromState(ctx context.Context, stateID bson.ObjectID, deckID bson.ObjectID) error {
	filter := bson.M{"_id": stateID}
	update := bson.M{
		"$pull": bson.M{
			"decks": bson.M{"id": deckID},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to remove deck: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("game state with ID %s not found", stateID.Hex())
	}

	r.logger.Infof("removed deck %s from game state %s", deckID, stateID.Hex())
	return nil
}

func (r *GameStateRepository) UpdateDeckInState(ctx context.Context, stateID bson.ObjectID, deckID bson.ObjectID, updateData *deckboard_models.Deck) error {
	filter := bson.M{
		"_id":      stateID,
		"decks.id": deckID,
	}

	update := bson.M{
		"$set": bson.M{
			"decks.$.metadata": updateData.Metadata,
			"decks.$.cards":    updateData.Cards,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to update deck: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("deck %s not found in game state %s", deckID, stateID.Hex())
	}

	r.logger.Infof("updated deck %s in game state %s", deckID, stateID.Hex())
	return nil
}

func (r *GameStateRepository) RemoveCardFromDeck(ctx context.Context, stateID, deckID, cardID bson.ObjectID) error {
	filter := bson.M{
		"_id":      stateID,
		"decks.id": deckID,
	}

	update := bson.M{
		"$pull": bson.M{
			"decks.$.cards": bson.M{"id": cardID},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to remove card from deck: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("game state %s or deck %s not found", stateID.Hex(), deckID)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("card %s not found in deck %s", cardID, deckID)
	}

	r.logger.Infof("removed card %s from deck %s in game state %s", cardID, deckID, stateID.Hex())
	return nil
}
