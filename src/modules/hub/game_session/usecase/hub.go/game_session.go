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

func NewGameSession(id string, game game_session_contracts.Game, gameName string) *GameSession {
	return &GameSession{
		ID:       id,
		Game:     game,
		GameName: gameName,
		Players:  make(map[string]*game_session_entity.Connection), // Инициализация карты
	}
}

func (s *GameSession) HandleMessage(userID string, data game_session_contracts.Action) {
	s.Game.HandleAction(userID, data)
}

func (s *GameSession) AddPlayer(conn *game_session_entity.Connection, isHost bool, mainColor, highlightColor string) {
	s.playersMU.Lock()
	s.Players[conn.UserID] = conn
	s.playersMU.Unlock()
	s.Game.AddPlayer(conn.UserID, isHost, mainColor, highlightColor)
}

func (s *GameSession) RemovePlayer(userID string) {
	s.playersMU.Lock()
	delete(s.Players, userID)
	s.playersMU.Unlock()
	s.Game.RemovePlayer(userID)
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

func (s *GameSession) BroadcastToAllExcept(excludeUserID string, message any) {
	s.playersMU.RLock()
	defer s.playersMU.RUnlock()

	if len(s.Players) <= 1 {
		return
	}

	uuids := make([]string, 0, len(s.Players)-1)
	for userID, p := range s.Players {
		if userID != excludeUserID {
			uuids = append(uuids, p.Kws.UUID)
		}
	}

	if len(uuids) == 0 {
		return
	}

	rawMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Ошибка сериализации broadcast-сообщения: %v", err)
		return
	}

	socketio.EmitToList(uuids, rawMessage)
}
