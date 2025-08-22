package game_session_hub

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/gofiber/contrib/socketio"
	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
	game_session_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/entity"
)

type GameSession struct {
	ID        string
	Game      game_session_contracts.Game
	GameName  string
	Players   map[string]*game_session_entity.Connection
	playersMU sync.RWMutex
}

func (s *GameSession) HandleMessage(userID string, data game_session_contracts.Action) {
	s.Game.HandleAction(userID, data)
}

func (s *GameSession) AddPlayer(conn *game_session_entity.Connection, isHost bool) {
	s.playersMU.Lock()
	s.Players[conn.UserID] = conn
	s.playersMU.Unlock()
}

func (s *GameSession) RemovePlayer(userID string) {
	s.playersMU.Lock()
	delete(s.Players, userID)
	s.playersMU.Unlock()
}

func (s *GameSession) SendToAll(message any) {
	s.playersMU.RLock()
	defer s.playersMU.RUnlock()

	if len(s.Players) == 0 {
		return
	}

	uuids := make([]string, 0, len(s.Players))
	for _, p := range s.Players {
		uuids = append(uuids, p.Kws.UUID)
	}

	rawMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Ошибка сериализации broadcast-сообщения: %v", err)
		return
	}

	socketio.EmitToList(uuids, rawMessage)
}

func (s *GameSession) SendToPlayer(userID string, message any) error {
	s.playersMU.RLock()
	defer s.playersMU.RUnlock()

	player, ok := s.Players[userID]
	if !ok {
		return errors.New("игрок не найден в сессии")
	}

	rawMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return socketio.EmitTo(player.Kws.UUID, rawMessage)
}

// func (h *Hub) CreateHub(ctx context.Context, hub *game_session_entity.Hub, gameName string) error {
// 	g, err := game_session_registry.NewGame(gameName)
// 	if err != nil {
// 		return err
// 	}
// 	hubb := game_session_entity.Hub{
// 		ID:      hub.ID,
// 		Game:    g,
// 		Players: make(map[string]*game_session_entity.Connection),
// 	}
// 	hubb.Game.Initialize(hub.ID)
// }
//
// func (h *Hub) broadcast(message any, recipients []string) error {
// 	h.hubMU.RLock()
// 	defer h.hubMU.RUnlock()
// 	hub, ok := h.hub[hubID]
// 	if !ok {
// 		return errors.New("hub not found")
// 	}
// 	hub.Game.HandleAction(message, recipients)
// 	return nil
// }
//
// func (h *Hub) JoinHub(ctx context.Context, conn *game_session_entity.Connection, hubID string) error {
// 	h.hubMU.Lock()
// 	defer h.hubMU.Unlock()
// 	hub, ok := h.hub[hubID]
// 	if !ok {
// 		return errors.New("hub not found")
// 	}
// 	hub.AddPlayer(conn.UserID, conn.Kws, conn.ISHost)
// 	return nil
// }
//
// func (h *Hub) LeaveHub(ctx context.Context, conn *game_session_entity.Connection, hubID string) error {
// 	h.hubMU.Lock()
// 	defer h.hubMU.Unlock()
// 	hub, ok := h.hub[hubID]
// 	if !ok {
// 		return nil
// 	}
// 	hub.RemovePlayer(conn.UserID)
// 	return nil
// }
