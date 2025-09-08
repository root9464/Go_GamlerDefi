package deckboard_models

// DICE
type ChangeDice struct {
	DiceCount   int `json:"dice_count"`
	FacesNumber int `json:"faces_number"`
}

type RollDiceResponse struct {
	ClientID string     `json:"client_id"`
	Dices    DiceResult `json:"dices"`
}

type DiceResult struct {
	Faces  int   `json:"faces"`
	Values []int `json:"value"`
}

// DECK
type GiveDeckForSelection struct {
	DeckID   string `json:"deck_id"`
	PlayerID string `json:"player_id"`
}

type PromptSelectCard struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	BackImageURL string   `json:"back_image_url"`
	Cards        []string `json:"cards"`
}

type SelectCard struct {
	DeckID string `json:"deck_id"`
	CardID string `json:"card_id"`
}

type GotDeck struct {
	Deck Deck `json:"deck"`
}

// TOKEN
type MoveToken struct {
	PlayerID string         `json:"player_id"`
	Position PlayerPosition `json:"position"`
}

// CARD
type CardReveal struct {
	Card Card `json:"card"`
}

type ShowPlayerHand struct {
	PlayerID string `json:"player_id"`
}
