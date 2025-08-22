package game_session_models

import "go.mongodb.org/mongo-driver/v2/bson"

type GameSession struct {
	ID     bson.ObjectID `bson:"_id,omitempty"`
	HostID string        `bson:"host_id"`
	GameID string        `bson:"game_id"`
}
