package deckboard_service

import (
	"errors"
	"fmt"
	"maps"
	"sync"

	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
	"github.com/shopspring/decimal"
)

type PlayerManager struct {
	players map[string]*deckboard_models.Player
	mu      sync.RWMutex
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		players: make(map[string]*deckboard_models.Player),
	}
}

func (pm *PlayerManager) SetMetadataValue(userID, key string, value any) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.players[userID]
	if !ok {
		return errors.New("player not found")
	}

	if player.Metadata == nil {
		player.Metadata = make(map[string]any)
	}

	player.Metadata[key] = value
	return nil
}

func (pm *PlayerManager) UpdateMetadata(userID string, data map[string]any) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.players[userID]
	if !ok {
		return errors.New("player not found")
	}

	if player.Metadata == nil {
		player.Metadata = make(map[string]any)
	}

	maps.Copy(player.Metadata, data)
	return nil
}

func (pm *PlayerManager) IncrementMetadataValue(userID, key string, amount decimal.Decimal) (decimal.Decimal, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.players[userID]
	if !ok {
		return decimal.Zero, errors.New("player not found")
	}

	if player.Metadata == nil {
		player.Metadata = make(map[string]any)
	}

	var currentValue decimal.Decimal
	if val, exists := player.Metadata[key]; exists {
		switch v := val.(type) {
		case float64:
			currentValue = decimal.NewFromFloat32(float32(v))
		case int:
			currentValue = decimal.NewFromInt(int64(v))
		case int32:
			currentValue = decimal.NewFromInt(int64(v))
		case int64:
			currentValue = decimal.NewFromInt(v)
		default:
			return decimal.Zero, fmt.Errorf("metadata key '%s' is not a number", key)
		}
	}

	newValue := decimal.Avg(currentValue, amount)
	player.Metadata[key] = newValue

	return newValue, nil
}

func (pm *PlayerManager) AddPlayer(id string, isHost bool, mainColor, highlightColor string) (*deckboard_models.Player, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if player, exist := pm.players[id]; exist {
		player.IsActive = true
		return player, nil
	}

	player := &deckboard_models.Player{
		ID:       id,
		Position: deckboard_models.PlayerPosition{X: decimal.NewFromFloat32(0), Y: decimal.NewFromFloat32(0)},
		Hand:     []deckboard_models.PlayerHand{},
		Metadata: make(map[string]any),
		Token: deckboard_models.Token{
			MainColor:      mainColor,
			HighlightColor: highlightColor,
		},
		IsHost:   isHost,
		IsActive: true,
	}

	pm.players[id] = player
	return player, nil
}

func (pm *PlayerManager) RemovePlayer(id string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, exist := pm.players[id]
	if exist {
		player.IsActive = false
	}

	return exist
}

func (pm *PlayerManager) GetPlayer(id string) (*deckboard_models.Player, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	player, ok := pm.players[id]
	if !ok {
		return nil, errors.New("player not found")
	}
	return player, nil
}

func (pm *PlayerManager) UpdatePosition(id string, x, y decimal.Decimal) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.players[id]
	if !ok {
		return errors.New("player not found")
	}

	player.Position.X = x
	player.Position.Y = y
	return nil
}

func (pm *PlayerManager) GiveCard(id, deckID string, card deckboard_models.Card) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.players[id]
	if !ok {
		return errors.New("player not found")
	}

	// Найти или создать руку для данной колоды
	for i := range player.Hand {
		if player.Hand[i].DeckID == deckID {
			player.Hand[i].Cards = append(player.Hand[i].Cards, card)
			return nil
		}
	}

	// Создать новую руку для колоды
	player.Hand = append(player.Hand, deckboard_models.PlayerHand{
		DeckID: deckID,
		Cards:  []deckboard_models.Card{card},
	})
	return nil
}

func (pm *PlayerManager) CollectCard(id, deckID, cardID string) (*deckboard_models.Card, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.players[id]
	if !ok {
		return nil, errors.New("player not found")
	}

	// Найти руку для данной колоды
	for i, hand := range player.Hand {
		if hand.DeckID == deckID {
			// Найти индекс карты
			cardIndex := -1
			for j, card := range hand.Cards {
				if card.ID == cardID {
					cardIndex = j
					break
				}
			}

			if cardIndex == -1 {
				return nil, errors.New("card not found in hand")
			}

			card := player.Hand[i].Cards[cardIndex]
			player.Hand[i].Cards = append(hand.Cards[:cardIndex], hand.Cards[cardIndex+1:]...)

			if len(player.Hand[i].Cards) == 0 {
				player.Hand = append(player.Hand[:i], player.Hand[i+1:]...)
			}

			return &card, nil
		}
	}

	return nil, errors.New("deck not found in player's hand")
}

func (pm *PlayerManager) GetAllPlayersState() []*deckboard_models.Player {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	playersState := make([]*deckboard_models.Player, 0, len(pm.players))
	for _, player := range pm.players {
		playerCopy := *player
		if player.IsActive {
			playersState = append(playersState, &playerCopy)
		}
	}
	return playersState
}
