package deckboard_models

import (
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Game struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Players     []Player `json:"players"`
	BoardURL    string   `json:"board_url"`
	Decks       []Deck   `json:"decks"`
}

type Player struct {
	ID       string         `json:"id"`
	Position PlayerPosition `bson:"position" json:"position"`
	Token    Token          `bson:"token" json:"token"`
	Hand     []Card         `json:"hand"`

	IsActive bool `bson:"is_active" json:"is_active"`

	Metadata bson.M `bson:"metadata" json:"metadata"`
}

type Token struct {
	MainColor      string `bson:"main_color" json:"main_color"`
	HighlightColor string `bson:"highlight_color" json:"highlight_color"`
}

type PlayerPosition struct {
	X decimal.Decimal `bson:"x" json:"x"`
	Y decimal.Decimal `bson:"y" json:"y"`
}

type Deck struct {
	ID bson.ObjectID `bson:"_id" json:"id"`
	// Name         string `json:"name" mapstructure:"name"`
	// BackImageURL string `json:"back_image_url" mapstructure:"back_image_url"`
	Cards    []Card `bson:"cards" json:"cards"`
	Metadata bson.M `bson:"metadata" json:"metadata"`
}

type Card struct {
	ID bson.ObjectID `bson:"_id" json:"id"`
	// Title       string `json:"title" `
	// Description string `json:"description" mapstructure:"description"`
	// Category    string `json:"category" mapstructure:"category"`
	// ImageURL    string `json:"image_url" mapstructure:"image_url"`
	// Task        string `json:"task" mapstructure:"task"`
	Metadata bson.M `bson:"metadata" json:"metadata"`
}
