package game_config_repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type GameConfigModel struct {
	ID string `bson:"_id"` // название самой игры

	Title       string `bson:"title"`
	Description string `bson:"description"`

	Settings primitive.M `bson:"settings"`
}
