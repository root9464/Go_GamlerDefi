package game_session_contracts

import game_config_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/repository"

type Action struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type Broadcast func(message any)

type Game interface {
	Initialize(
		hostID string,
		sendToAll func(message any),
		sendToPlayer func(userID string, message any) error,
		broadcastToAllExcept func(excludeUserID string, message any),
	)
	AddPlayer(userID string, isHost bool, mainColor, highlightColor string)
	RemovePlayer(userID string)
	HandleAction(clientID string, action Action)
}

type SettingsProvider interface {
	GetSettingsModel() any
}

type GameFactory func(repo game_config_repository.GameConfigRepository) (Game, error)
