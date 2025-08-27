package game_session_models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SessionStatus string

const (
	StatusScheduled SessionStatus = "scheduled"
	StatusActive    SessionStatus = "active"
	StatusFinished  SessionStatus = "finished"
)

type Player struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

type GameSession struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	HostID       string        `bson:"host_id"`
	GameName     string        `bson:"game_name"`    // Важно хранить имя игры для реестра
	Participants []Player      `bson:"participants"` // ID игроков, которые "купили билет"
	Status       SessionStatus `bson:"status"`
	TimeToStart  time.Time     `bson:"time_to_start"`
	Price        float64       `bson:"price"`
}
