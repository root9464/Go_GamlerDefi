package game_session_registry

import (
	"fmt"
	"sync"

	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
	game_config_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/repository"
)

type RegistredGame struct {
	Factory          game_session_contracts.GameFactory
	SettingsProvider game_session_contracts.SettingsProvider
}

var (
	registeredGames = make(map[string]RegistredGame)
	mu              sync.RWMutex
)

func RegisterGame(name string, factory game_session_contracts.GameFactory, settingsProvider game_session_contracts.SettingsProvider) {
	mu.Lock()
	defer mu.Unlock()

	registeredGames[name] = RegistredGame{
		Factory:          factory,
		SettingsProvider: settingsProvider,
	}

	fmt.Printf("INFO: Игра '%s' успешно зарегистрирована в системе.\n", name)
}

func NewGame(name string, repo game_config_repository.GameConfigRepository) (game_session_contracts.Game, error) {
	mu.RLock()
	game, ok := registeredGames[name]
	mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("game '%s' not found", name)
	}

	return game.Factory(repo)
}

func IsGameRegistered(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := registeredGames[name]
	return ok
}

func GetGameSettingsModel(name string) (any, bool) {
	mu.RLock()
	defer mu.RUnlock()
	game, ok := registeredGames[name]
	if !ok {
		return nil, false
	}
	return game.SettingsProvider.GetSettingsModel(), true
}
