package sales_courage_entity

import "time"

type Game struct {
	ID          string
	Title       string
	Description string

	BoardURL string
	Decks    []Deck

	Players []Player

	GameSettings GameSettings
}

type Player struct {
	ID       string
	Name     string
	Position any
	Coins    int
	Hand     []Card
}

type GameSettings struct {
	MaxPlayers int
	MinPLayers int

	MaxDiceCount int

	TimeLimit time.Duration
}

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
