package deckboard_service

import (
	"context"
	"fmt"

	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IDeckRepository interface {
	AddDecksToState(ctx context.Context, decks []*deckboard_models.Deck, stateID bson.ObjectID) error
	GetDeckFromState(ctx context.Context, stateID bson.ObjectID, deckID bson.ObjectID) (*deckboard_models.Deck, error)
	GetAllDecksFromState(ctx context.Context, stateID bson.ObjectID) ([]deckboard_models.Deck, error)
	RemoveDeckFromState(ctx context.Context, stateID bson.ObjectID, deckID bson.ObjectID) error
	UpdateDeckInState(ctx context.Context, stateID bson.ObjectID, deckID bson.ObjectID, updateData *deckboard_models.Deck) error

	RemoveCardFromDeck(ctx context.Context, stateID, deckID, cardID bson.ObjectID) error
}

type IDeckManager interface {
	GetAllDecks(ctx context.Context) ([]deckboard_models.Deck, error)
	GetDeck(ctx context.Context, deckID bson.ObjectID) (*deckboard_models.Deck, error)
	DrawCard(ctx context.Context, deckID, cardID bson.ObjectID) (*deckboard_models.Card, error)
}

type DeckManager struct {
	stateID    bson.ObjectID
	logger     *logger.Logger
	repository IDeckRepository
}

func NewDeckManager(
	logger *logger.Logger,
	repository IDeckRepository,
	stateID bson.ObjectID,
) IDeckManager {
	return &DeckManager{
		logger:     logger,
		repository: repository,
		stateID:    stateID,
	}
}

func (dm *DeckManager) GetAllDecks(ctx context.Context) ([]deckboard_models.Deck, error) {
	decks, err := dm.repository.GetAllDecksFromState(ctx, dm.stateID)
	if err != nil {
		return nil, err
	}

	return decks, nil
}

func (dm *DeckManager) GetDeck(ctx context.Context, deckID bson.ObjectID) (*deckboard_models.Deck, error) {
	deck, err := dm.repository.GetDeckFromState(ctx, dm.stateID, deckID)
	if err != nil {
		return nil, err
	}

	return deck, nil
}

func (dm *DeckManager) DrawCard(ctx context.Context, deckID, cardID bson.ObjectID) (*deckboard_models.Card, error) {
	deck, err := dm.repository.GetDeckFromState(ctx, dm.stateID, deckID)
	if err != nil {
		return nil, err
	}

	for _, card := range deck.Cards {
		if card.ID == cardID {
			if err := dm.repository.RemoveCardFromDeck(ctx, dm.stateID, deckID, cardID); err != nil {
				return nil, err
			}
			return &card, nil
		}
	}

	return nil, fmt.Errorf("card not found")
}
