package game_session_contracts

type Action struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type Broadcast func(message any)

type Game interface {
	Initialize(
		sessionID string,
		sendToAll func(message any),
		sendToPlayer func(userID string, message any) error,
	)
	HandleAction(clientID string, action Action)
}
