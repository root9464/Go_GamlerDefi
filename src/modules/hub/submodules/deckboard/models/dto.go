package deckboard_models

type RollDice struct {
	DiceCount   int `json:"dice_count"`
	FacesNumber int `json:"faces_number"`
}

type RollDiceResponse struct {
	ClientID string `json:"client_id"`
	Dices    []DiceResult
}

type DiceResult struct {
	Faces int `json:"faces"`
	Value int `json:"value"`
}
