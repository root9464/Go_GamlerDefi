package game_config_model

type GameConfigModel struct {
	ID string `bson:"_id"` // название самой игры

	Title       string `bson:"title"`
	Description string `bson:"description"`

	Settings any `bson:"settings"`
}
