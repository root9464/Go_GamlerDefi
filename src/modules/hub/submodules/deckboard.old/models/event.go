package deckboard_models

type FullGameStatePayload struct {
	SessionID string    `json:"sessionId"`
	GameTitle string    `json:"gameTitle"`
	Players   []*Player `json:"players"`
}

type PlayerJoinedPayload struct {
	Player *Player `json:"player"`
}

type PlayerLeftPayload struct {
	PlayerID string `json:"playerId"`
}

type TokenMovedPayload struct {
	PlayerID string         `json:"playerId"`
	Position PlayerPosition `json:"position"`
}
