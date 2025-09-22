package deckboard_dto

import "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"

// DICE
type ChangeDiceSettings struct {
	DiceCount int `json:"dice_count"`
	DiceFaces int `json:"dice_faces"`
}

// PLAYER
type AddPlayer struct {
	UserID   string         `json:"user_id"`
	Metadata map[string]any `json:"metadata"`
}

type MoveToken struct {
	Position deckboard_models.PlayerPosition `json:"position"`
}
