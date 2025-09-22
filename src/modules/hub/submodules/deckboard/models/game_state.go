package deckboard_models

import "go.mongodb.org/mongo-driver/v2/bson"

type GameState struct {
	ID bson.ObjectID `bson:"_id" json:"id"`

	Players []Player `bson:"players" json:"players"`
	Decks   []Deck   `bson:"decks" json:"decks"`

	DiceState DiceState `bson:"dice_state" json:"dice_state"`

	HostID string `bson:"host_id" json:"host_id"`
}

type DiceState struct {
	DiceCount int `bson:"dice_count" json:"dice_count"`
	DiceFaces int `bson:"dice_faces" json:"dice_faces"`
}

func NewGameState(
	decks []Deck,
	diceCount, diceFaces int,
	hostID string,
) *GameState {
	newDecks := make([]Deck, 0)
	newDecks = append(newDecks, decks...)
	return &GameState{
		Players: make([]Player, 0),
		Decks:   newDecks,
		DiceState: DiceState{
			DiceCount: diceCount,
			DiceFaces: diceFaces,
		},
		HostID: hostID,
	}
}
