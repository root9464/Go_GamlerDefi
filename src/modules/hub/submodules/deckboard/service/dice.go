package deckboard_service

import (
	"math/rand"
	"time"

	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
)

type DiceManager struct{}

func NewDiceManager() *DiceManager {
	return &DiceManager{}
}

func (dm *DiceManager) RollDices(diceCount int, facesNumber int) []deckboard_models.DiceResult {
	results := make([]deckboard_models.DiceResult, diceCount)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range results {
		results[i] = deckboard_models.DiceResult{
			Faces: facesNumber,
			Value: r.Intn(facesNumber) + 1,
		}
	}

	return results
}
