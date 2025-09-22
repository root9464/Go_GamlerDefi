package games_submodules_utils

import (
	"encoding/json"

	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
)

func EncodeMsg(action *game_session_contracts.Action, dest any) error {
	if action.Payload == nil {
		return nil
	}

	payloadBytes, err := json.Marshal(action.Payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(payloadBytes, dest)
}
