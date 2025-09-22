package deckboard_handlers

import (
	"github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/service"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type DeckboardHandlers struct {
	logger *logger.Logger

	SendToAll            func(message any)
	SendToPlayer         func(userID string, message any) error
	BroadcastToAllExcept func(excludeUserID string, message any)

	playerManager deckboard_service.IPlayerManager
	deckManager   deckboard_service.IDeckManager
	diceManager   deckboard_service.IDiceManager
}

type IDeckboardHandlers interface {
	iHttpHandlers
	iWsHandlers
}

func NewDeckboardHandlers(
	logger *logger.Logger,
) IDeckboardHandlers {
	return &DeckboardHandlers{
		logger: logger,
	}
}
