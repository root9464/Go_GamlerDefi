package game_session_entity

type GameSession struct {
	ID     string
	HostID string
	GameID string

	GameState GameState
}

type GameState struct{}
