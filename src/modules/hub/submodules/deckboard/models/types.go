package deckboard_models

type GameConfig struct {
	Title        string
	Description  string
	Decks        []Deck
	MaxDiceCount int
}

type DataProvider interface {
	GetGameConfig() GameConfig
}
