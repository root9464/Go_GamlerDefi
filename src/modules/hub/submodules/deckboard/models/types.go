package deckboard_models

type GameConfig struct {
	HostID      string
	Title       string
	Description string
	Decks       []Deck

	DiceCount int `json:"dice_count"`
	DiceFaces int `json:"dice_faces"`
}

type DataProvider interface {
	GetGameConfig() GameConfig
}
