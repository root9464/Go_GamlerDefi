package game_session_entity

import (
	"sync"

	"github.com/gofiber/contrib/socketio"
	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
)

type Connection struct {
	UserID string
	ISHost bool
	Kws    *socketio.Websocket
}

type Hub struct {
	ID        string                 // id сессии
	Players   map[string]*Connection // id userID => connection
	playersMU sync.RWMutex
	Game      game_session_contracts.Game
}

func (gs *Hub) AddPlayer(userID string, kws *socketio.Websocket, isHost bool) {
	gs.playersMU.Lock()
	defer gs.playersMU.Unlock()
	gs.Players[userID] = &Connection{
		UserID: userID,
		Kws:    kws,
		ISHost: isHost,
	}
}

func (gs *Hub) RemovePlayer(userID string) {
	gs.playersMU.Lock()
	defer gs.playersMU.Unlock()
	delete(gs.Players, userID)
}

func (gs *Hub) GetPlayer(userID string) *Connection {
	gs.playersMU.RLock()
	defer gs.playersMU.RUnlock()
	return gs.Players[userID]
}

func (gs *Hub) GetPlayers() map[string]*Connection {
	gs.playersMU.RLock()
	defer gs.playersMU.RUnlock()
	return gs.Players
}

func (gs *Hub) IsHost(userID string) bool {
	gs.playersMU.RLock()
	defer gs.playersMU.RUnlock()
	return gs.Players[userID].ISHost
}

func (gs *Hub) IsPlayer(userID string) bool {
	gs.playersMU.RLock()
	defer gs.playersMU.RUnlock()
	_, ok := gs.Players[userID]
	return ok
}
