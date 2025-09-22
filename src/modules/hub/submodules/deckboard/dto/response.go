package deckboard_dto

import "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"

// DICE
type RollDiceResult struct {
	Faces  int   `json:"faces"`
	Values []int `json:"value"`
}

// PLAYER
type TokenMoved struct {
	PlayerID string                          `json:"playerId"`
	Position deckboard_models.PlayerPosition `json:"position"`
}
