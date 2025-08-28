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

func (pm *PlayerManager) IncrementMetadataValue(userID, key string, amount float64) (float64, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.players[userID]
	if !ok {
		return 0, errors.New("player not found")
	}

	if player.Metadata == nil {
		player.Metadata = make(map[string]any)
	}

	var currentValue float64
	if val, exists := player.Metadata[key]; exists {
		switch v := val.(type) {
		case float64:
			currentValue = v
		case int:
			currentValue = float64(v)
		case int32:
			currentValue = float64(v)
		case int64:
			currentValue = float64(v)
		default:
			return 0, fmt.Errorf("metadata key '%s' is not a number", key)
		}
	}

	newValue := currentValue + amount
	player.Metadata[key] = newValue

	return newValue, nil
}

func (pm *PlayerManager) AddPlayer(id string, isHost bool) (*deckboard_models.Player, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exist := pm.players[id]; exist {
		return nil, errors.New("player already exists")
	}

	player := &deckboard_models.Player{
		ID:       id,
		Position: deckboard_models.PlayerPosition{X: decimal.NewFromFloat32(0), Y: decimal.NewFromFloat32(0)},
		Hand:     []deckboard_models.Card{},
	}

	pm.players[id] = player
	return player, nil
}

func (pm *PlayerManager) RemovePlayer(id string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	_, exist := pm.players[id]
	if exist {
		delete(pm.players, id)
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

func (pm *PlayerManager) GiveCard(id string, card deckboard_models.Card) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.players[id]
	if !ok {
		return errors.New("player not found")
	}

	player.Hand = append(player.Hand, card)
	return nil
}

func (pm *PlayerManager) GetAllPlayersState() []*deckboard_models.Player {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	playersState := make([]*deckboard_models.Player, 0, len(pm.players))
	for _, player := range pm.players {
		playerCopy := *player
		playersState = append(playersState, &playerCopy)
	}
	return playersState
}
