package game_session_entity

import (
	"time"
)

type SessionStatus string

const (
	StatusScheduled SessionStatus = "scheduled"
	StatusActive    SessionStatus = "active"
	StatusFinished  SessionStatus = "finished"
)

type Player struct {
	ID   string
	Name string
}

type GameSession struct {
	ID           string
	HostID       string
	GameName     string
	Participants []Player
	Status       SessionStatus
	TimeToStart  time.Time
	Price        float64
}
