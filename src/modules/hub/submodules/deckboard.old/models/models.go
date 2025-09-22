package deckboard_models

import (
	"time"

	"github.com/shopspring/decimal"
)

// type Deck struct {
// 	ID           string `json:"id"`
// 	Name         string `json:"name"`
// 	BackImageURL string `json:"back_image_url"`
// 	Cards        []Card `json:"cards"`
// }
//
// type Card struct {
// 	ID          string `json:"id"`
// 	Title       string `json:"title"`
// 	Description string `json:"description"`
// 	Category    string `json:"category"`
// 	ImageURL    string `json:"image_url"`
// 	Task        string `json:"task"`
// }

type Deck struct {
	ID           string `json:"id" mapstructure:"id"`
	Name         string `json:"name" mapstructure:"name"`
	BackImageURL string `json:"back_image_url" mapstructure:"back_image_url"`
	Cards        []Card `json:"cards" mapstructure:"cards"`
}

// Card - модель карты
type Card struct {
	ID          string `json:"id" mapstructure:"id"`
	Title       string `json:"title" mapstructure:"title"`
	Description string `json:"description" mapstructure:"description"`
	Category    string `json:"category" mapstructure:"category"`
	ImageURL    string `json:"image_url" mapstructure:"image_url"`
	Task        string `json:"task" mapstructure:"task"`
}

type Game struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Players     []Player  `json:"players"`
	BoardURL    string    `json:"board_url"`
	Decks       []Deck    `json:"decks"`
	CreateAt    time.Time `json:"create_at"`
	StartedAt   time.Time `json:"started_at"`
	FinishedAt  time.Time `json:"finished_at"`
}

type Player struct {
	ID       string         `json:"id"`
	Position PlayerPosition `json:"position"`
	Hand     []PlayerHand   `json:"hand"`
	Token    Token          `bson:"token" json:"token"`

	Metadata map[string]any `json:"metadata"`
	IsActive bool           `json:"is_active"`
}

type Token struct {
	MainColor      string `bson:"main_color" json:"main_color"`
	HighlightColor string `bson:"highlight_color" json:"highlight_color"`
}

type PlayerHand struct {
	DeckID string `json:"deck_id"`
	Cards  []Card `json:"cards"`
}

type PlayerPosition struct {
	X decimal.Decimal `json:"x"`
	Y decimal.Decimal `json:"y"`
}
