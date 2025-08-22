package game_session_registry

import (
	"fmt"
	"sync"

	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
)

var (
	gameFactories = make(map[string]func() game_session_contracts.Game)
	mu            sync.RWMutex
)

func RegisterGame(name string, factory func() game_session_contracts.Game) {
	// log.Fatalf("Name: %s, Factory: %v", name, factory)
	mu.Lock()
	defer mu.Unlock()
	gameFactories[name] = factory
}

func NewGame(name string) (game_session_contracts.Game, error) {
	mu.RLock()
	defer mu.RUnlock()
	factory, ok := gameFactories[name]
	if !ok {
		return nil, fmt.Errorf("игра '%s' не найдена", name)
	}
	return factory(), nil
}
