package deckboard_models

import "time"

type Deck struct {
	ID           string
	Name         string
	BackImageURL string

	Cards []Card
}

type Card struct {
	ID          string
	Title       string
	Description string
	Category    string
	ImageURL    string
	Task        string
}

type Game struct {
	ID          string
	Title       string
	Description string
	Players     []Player

	BoardURL string
	Decks    []Deck

	GameSettings GameSettings

	CreateAt   time.Time
	StartedAt  time.Time
	FinishedAt time.Time
}

type GameSettings struct {
	MaxPlayers int
	MinPlayers int

	MinDuration time.Duration
	MaxDuration time.Duration

	MaxDiceCount int
	TimeLimit    time.Duration
}

type GameState struct {
	DiceCount int
}

type Player struct {
	ID     string
	Name   string
	Avatar string

	Score    int
	Position any

	PersonalCards []Card
	Metadata      map[string]any

	IsActive   bool
	IsFinished bool
}
