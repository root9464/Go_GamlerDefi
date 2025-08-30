package game_config_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateGameConfigDTO struct {
	GameName    string      `json:"game_name"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Settings    primitive.M `json:"settings"`
}
