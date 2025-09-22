package deckboard_service

import (
	"math/rand"
	"time"

	deckboard_dto "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/dto"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type IDiceManager interface {
	RollDices() *deckboard_dto.RollDiceResult
	ChangeDiceSettings(data *deckboard_dto.ChangeDiceSettings) error
}

type DiceManager struct {
	logger *logger.Logger

	diceCount int
	diceFaces int
}

func NewDiceManager(logger *logger.Logger) IDiceManager {
	return &DiceManager{
		logger: logger,
	}
}

func (dm *DiceManager) RollDices() *deckboard_dto.RollDiceResult {
	dm.logger.Infof("rolling %d dice with %d faces", dm.diceCount, dm.diceFaces)
	result := deckboard_dto.RollDiceResult{
		Faces:  dm.diceFaces,
		Values: make([]int, dm.diceCount),
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range result.Values {
		result.Values[i] = r.Intn(dm.diceFaces) + 1
	}

	return &result
}

func (dm *DiceManager) ChangeDiceSettings(data *deckboard_dto.ChangeDiceSettings) error {
	if data.DiceCount > 0 {
		dm.diceCount = data.DiceCount
	}

	if data.DiceFaces > 0 {
		dm.diceFaces = data.DiceFaces
	}

	dm.logger.Infof("change dice successful. New count: %d; New faces: %d", dm.diceCount, dm.diceFaces)
	return nil
}
