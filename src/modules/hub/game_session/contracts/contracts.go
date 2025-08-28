package game_session_contracts

type Action struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type Broadcast func(message any)

type Game interface {
	Initialize(
		sendToAll func(message any),
		sendToPlayer func(userID string, message any) error,
		broadcastToAllExcept func(excludeUserID string, message any),
	)
	AddPlayer(userID string, isHost bool)
	RemovePlayer(userID string)
	HandleAction(clientID string, action Action)
}
