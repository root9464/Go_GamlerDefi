package deckboard_service

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
)

type DeckManager struct {
	decks map[string]*deckboard_models.Deck
	mu    sync.RWMutex
}

func NewDeckManager(decks []deckboard_models.Deck) *DeckManager {
	decksMap := make(map[string]*deckboard_models.Deck)
	for _, deck := range decks {
		decksMap[deck.ID] = &deck
	}

	return &DeckManager{
		decks: decksMap,
	}
}

func (dm *DeckManager) GetDecks() map[string]*deckboard_models.Deck {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.decks
}

func (dm *DeckManager) GetDeck(deckID string) (*deckboard_models.Deck, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	deck, exists := dm.decks[deckID]
	if !exists {
		return nil, errors.New("Deck not found")
	}
	return deck, nil
}

func (dm *DeckManager) ReturnCard(deckID string, card deckboard_models.Card) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	deck, exists := dm.decks[deckID]
	if !exists {
		return errors.New("deck not found")
	}

	deck.Cards = append(deck.Cards, card)
	dm.decks[deckID] = deck

	return nil
}

func (dm *DeckManager) DrawCard(deckID string, cardID string) (*deckboard_models.Card, error) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	deck, exists := dm.decks[deckID]
	if !exists || len(deck.Cards) == 0 {
		return nil, errors.New("deck not found or empty")
	}

	for i, card := range deck.Cards {
		if card.ID == cardID {
			deck.Cards = append(deck.Cards[:i], deck.Cards[i+1:]...)
			return &card, nil
		}
	}
	return nil, errors.New("card not found")
}

func (dm *DeckManager) ShuffleDeck(deckID string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	deck, exists := dm.decks[deckID]
	if !exists {
		return errors.New("deck not found")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	shuffled := make([]deckboard_models.Card, len(deck.Cards))
	copy(shuffled, deck.Cards)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	deck.Cards = shuffled
	dm.decks[deckID] = deck

	return nil
}
