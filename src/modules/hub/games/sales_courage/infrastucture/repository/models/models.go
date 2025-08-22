package sales_courage_models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Game struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	BoardURL    string        `bson:"board_url"`
	Decks       []Deck        `bson:"decks"`
	// Players      []Player      `bson:"players"`
	GameSettings GameSettings `bson:"game_settings"`
	CreatedAt    time.Time    `bson:"created_at"`
	UpdatedAt    time.Time    `bson:"updated_at"`
	IsActive     bool         `bson:"is_active"`
}

// type Player struct {
// 	ID       bson.ObjectID `bson:"_id,omitempty"`
// 	Name     string        `bson:"name"`
// 	Position any           `bson:"position"`
// 	Hand     []Card        `bson:"hand"`
// 	Coins    int           `bson:"coins"`
// 	JoinedAt time.Time     `bson:"joined_at"`
// }

type GameSettings struct {
	MaxPlayers   int           `bson:"max_players"`
	MinPlayers   int           `bson:"min_players"`
	MaxDiceCount int           `bson:"max_dice_count"`
	TimeLimit    time.Duration `bson:"time_limit"`
}

type Deck struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	Name         string        `bson:"name"`
	BackImageURL string        `bson:"back_image_url"`
	Cards        []Card        `bson:"cards"`
	IsActive     bool          `bson:"is_active"`
	CreatedAt    time.Time     `bson:"created_at"`
}

type Card struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	Category    string        `bson:"category"`
	ImageURL    string        `bson:"image_url"`
	Task        string        `bson:"task"`
	CreatedAt   time.Time     `bson:"created_at"`
}
