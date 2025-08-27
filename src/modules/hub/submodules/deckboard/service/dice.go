package deckboard_service

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2/log"
	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
)

type DiceManager struct{}

func NewDiceManager() *DiceManager {
	return &DiceManager{}
}

func (dm *DiceManager) RollDices(diceCount int, facesNumber int) deckboard_models.DiceResult {
	log.Infof("Rolling %d dice with %d faces", diceCount, facesNumber)
	results := deckboard_models.DiceResult{
		Faces:  facesNumber,
		Values: make([]int, diceCount),
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range results.Values {
		results.Values[i] = r.Intn(facesNumber) + 1
	}

	return results
}
