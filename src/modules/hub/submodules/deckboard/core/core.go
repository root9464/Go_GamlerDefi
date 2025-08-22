package deckboard_core

import (
	"encoding/json"
	"fmt"

	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
	deckboard_service "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/service"
)

type DiceRollEvent struct {
	EventType string                            `json:"eventType"` // Хорошая практика для фронтенда
	Message   string                            `json:"message"`   // Человекочитаемое сообщение
	Details   deckboard_models.RollDiceResponse `json:"details"`   // Структурированные данные броска
}

type Template struct {
	State    *deckboard_models.GameState
	Provider deckboard_models.DataProvider

	SendToAll    func(message any)
	SendToPlayer func(userID string, message any) error

	diceManager *deckboard_service.DiceManager
	deckManager *deckboard_service.DeckManager
}

func (t *Template) Inizialize(sendToAll func(message any), sendToPlayer func(userID string, message any) error, provider deckboard_models.DataProvider) {
	t.Provider = provider

	t.SendToAll = sendToAll
	t.SendToPlayer = sendToPlayer

	config := t.Provider.GetGameConfig()

	t.State = &deckboard_models.GameState{}

	t.diceManager = deckboard_service.NewDiceManager()
	t.deckManager = deckboard_service.NewDeckManager(config.Decks)
}

func (t *Template) HandleAction(clientID string, action *game_session_contracts.Action) {
	switch action.Type {
	case "roll_dice":
		data := new(deckboard_models.RollDice)
		if err := t.encodeMsg(action, data); err != nil {
			return
		}

		diceRes := t.diceManager.RollDices(data.DiceCount, data.FacesNumber)

		payload := DiceRollEvent{
			EventType: "dice_roll",
			Message:   fmt.Sprintf("пользователь %s бросил %d кубиков", clientID, data.DiceCount),
			Details: deckboard_models.RollDiceResponse{
				ClientID: clientID,
				Dices:    diceRes,
			},
		}

		t.SendToAll(payload)

		// case "draw_card":
		// 	data := new(deckboard_models.DrawCard)
		// 	if err := t.encodeMsg(action, data); err != nil {
		// 		return
		// 	}
		//
		// 	card, err := t.deckManager.DrawCard(data.DeckID, data.CardID)
		// 	if err != nil {
		// 		return
		// 	}
		// 	res := deckboard_models.DrawCardResponse{
		// 		ClientID: clientID,
		// 		Card:     *card,
		// 	}
		// 	resByte, err := json.Marshal(res)
		// 	if err != nil {
		// 		return
		// 	}
		// 	t.Broadcast(resByte)

		// пользователь может показать карту всем игрокам
	}
}

func (t *Template) encodeMsg(action *game_session_contracts.Action, dest any) error {
	if action.Payload == nil {
		return nil
	}

	payloadBytes, err := json.Marshal(action.Payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(payloadBytes, dest)
}
