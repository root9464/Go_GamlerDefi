package deckboard_handlers

import (
	"context"

	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
	deckboard_dto "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/dto"
	games_submodules_utils "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/utils"
)

type iWsHandlers interface {
	HandleAction(ctx context.Context, clientID string, action *game_session_contracts.Action) error
}

func (h *DeckboardHandlers) HandleAction(ctx context.Context, clientID string, action *game_session_contracts.Action) error {
	switch action.Type {
	case "move_token":
		data := new(deckboard_dto.MoveToken)
		if err := games_submodules_utils.EncodeMsg(action, data); err != nil {
			return err
		}

		if err := h.playerManager.UpdatePlayerPosition(ctx, clientID, &data.Position); err != nil {
			return err
		}

		event := game_session_contracts.Action{
			Type: "token_moved",
			Payload: deckboard_dto.TokenMoved{
				PlayerID: clientID,
				Position: data.Position,
			},
		}

		h.BroadcastToAllExcept(clientID, event)
	}

	return nil
}
