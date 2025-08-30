package game_session_entity

import (
	"time"
)

type SessionStatus string

type GameSession struct {
	ID          string
	HostID      string
	GameName    string
	TimeToStart time.Time
}
