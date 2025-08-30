package game_session_models

import (
	"time"
)

type GameSession struct {
	ID          string    `bson:"_id,omitempty"`
	HostID      string    `bson:"host_id"`
	GameName    string    `bson:"game_name"` // Важно хранить имя игры для реестра
	TimeToStart time.Time `bson:"time_to_start"`
}
