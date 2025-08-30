package deckboard_core

import (
	"encoding/json"
	"fmt"
	"log"

	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
	deckboard_service "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/service"
)

type DiceRollEvent struct {
	EventType string                            `json:"eventType"`
	Message   string                            `json:"message"`
	Details   deckboard_models.RollDiceResponse `json:"details"`
}

type FullGameState struct {
	SessionID string                     `json:"sessionId"`
	GameTitle string                     `json:"gameTitle"`
	Players   []*deckboard_models.Player `json:"players"`
}

type Template struct {
	Provider deckboard_models.DataProvider
	config   deckboard_models.GameConfig

	HostID string

	SendToAll            func(message any)
	SendToPlayer         func(userID string, message any) error
	BroadcastToAllExcept func(excludeUserID string, message any)

	diceManager   *deckboard_service.DiceManager
	deckManager   *deckboard_service.DeckManager
	PlayerManager *deckboard_service.PlayerManager
}

func (t *Template) Inizialize(
	hostID string,

	sendToAll func(message any),
	sendToPlayer func(userID string, message any) error,
	broadcastToAllExcept func(excludeUserID string, message any),

	provider deckboard_models.DataProvider,
) {
	t.Provider = provider
	t.HostID = hostID

	t.SendToAll = sendToAll
	t.SendToPlayer = sendToPlayer
	t.BroadcastToAllExcept = broadcastToAllExcept

	t.config = t.Provider.GetGameConfig()

	t.diceManager = deckboard_service.NewDiceManager()
	t.deckManager = deckboard_service.NewDeckManager(t.config.Decks)
	t.PlayerManager = deckboard_service.NewPlayerManager()
}

func (t *Template) HandleHostAction(clientID string, action *game_session_contracts.Action, handler func()) {
	if clientID != t.HostID {
		t.SendToPlayer(clientID, "Это действие доступно только ведущему.")
		return
	}
	handler()
}

func (t *Template) AddPlayer(userID string, isHost bool) {
	newPlayer, err := t.PlayerManager.AddPlayer(userID, isHost)
	if err != nil {
		t.SendFullStateToPlayer(userID)
		return
	}

	t.SendFullStateToPlayer(userID)
	event := game_session_contracts.Action{
		Type:    "player_joined",
		Payload: deckboard_models.PlayerJoinedPayload{Player: newPlayer},
	}
	t.BroadcastToAllExcept(userID, event)
}

func (t *Template) RemovePlayer(userID string) {
	if t.PlayerManager.RemovePlayer(userID) {
		event := game_session_contracts.Action{
			Type:    "player_left",
			Payload: deckboard_models.PlayerLeftPayload{PlayerID: userID},
		}
		t.SendToAll(event)
	}
}

func (t *Template) HandleAction(clientID string, action *game_session_contracts.Action) {
	switch action.Type {
	case "move_token":
		data := new(deckboard_models.MoveToken)
		if err := t.EncodeMsg(action, data); err != nil {
			return
		}

		t.PlayerManager.UpdatePosition(clientID, data.Position.X, data.Position.Y)
		event := game_session_contracts.Action{
			Type: "token_moved",
			Payload: deckboard_models.TokenMovedPayload{
				PlayerID: clientID,
				Position: data.Position,
			},
		}

		log.Println(event)

		t.SendToAll(event)

	case "roll_dice":
		diceRes := t.diceManager.RollDices(t.config.DiceCount, t.config.DiceFaces)

		event := game_session_contracts.Action{
			Type: "dice_rolled",
			Payload: deckboard_models.RollDiceResponse{
				ClientID: clientID,
				Dices:    diceRes,
			},
		}
		t.SendToAll(event)

	case "select_card":
		data := new(deckboard_models.SelectCard)
		if err := t.EncodeMsg(action, data); err != nil {
			return
		}

		card, err := t.deckManager.DrawCard(data.DeckID, data.CardID)
		if err != nil {
			return
		}

		event := game_session_contracts.Action{
			Type:    "card_selected",
			Payload: card,
		}

		t.SendToPlayer(clientID, event)

	case "show_everyone_card":
		data := new(deckboard_models.SelectCard)
		if err := t.EncodeMsg(action, data); err != nil {
			return
		}

		player, err := t.PlayerManager.GetPlayer(clientID)
		if err != nil {
			return
		}

		for _, card := range player.Hand {
			if card.ID == data.CardID {
				event := game_session_contracts.Action{
					Type: "card_revealed",
					Payload: deckboard_models.CardReveal{
						Card: card,
					},
				}
				t.SendToAll(event)
				break
			}
		}

	// Host methods
	case "change_dice":
		t.HandleHostAction(clientID, action, func() {
			data := new(deckboard_models.ChangeDice)
			if err := t.EncodeMsg(action, data); err != nil {
				return
			}

			if data.DiceCount > 0 {
				t.config.DiceCount = data.DiceCount
			}

			if data.FacesNumber > 0 {
				t.config.DiceFaces = data.FacesNumber
			}

			payload := fmt.Sprintf("Dice changed to %d %d", t.config.DiceCount, t.config.DiceFaces)
			t.SendToPlayer(clientID, payload)
		})

	case "get_decks":
		t.HandleHostAction(clientID, action, func() {
			decks := t.deckManager.GetDecks()
			decksSlice := make([]deckboard_models.GotDeck, 0, len(decks))
			for _, deck := range decks {
				decksSlice = append(decksSlice, deckboard_models.GotDeck{
					Deck: deck,
				})
			}

			event := game_session_contracts.Action{
				Type:    "got_decks",
				Payload: decksSlice,
			}

			t.SendToPlayer(clientID, event)
		})

	case "give_deck_for_selection":
		t.HandleHostAction(clientID, action, func() {
			data := new(deckboard_models.GiveDeckForSelection)
			if err := t.EncodeMsg(action, data); err != nil {
				return
			}

			deck, err := t.deckManager.GetDeck(data.DeckID)
			if err != nil {
				return
			}

			cards := make([]string, len(deck.Cards))
			for _, card := range deck.Cards {
				cards = append(cards, card.ID)
			}
			event := game_session_contracts.Action{
				Type: "prompt_select_card",
				Payload: deckboard_models.PromptSelectCard{
					ID:           deck.ID,
					Name:         deck.Name,
					BackImageURL: deck.BackImageURL,
					Cards:        cards,
				},
			}

			hostPayload := fmt.Sprintf("Card drawn from deck %s", data.DeckID)

			t.SendToPlayer(data.PlayerID, event)
			t.SendToPlayer(clientID, hostPayload)
		})
	}
}

func (t *Template) SendFullStateToPlayer(userID string) {
	fullStatePayload := deckboard_models.FullGameStatePayload{
		GameTitle: t.Provider.GetGameConfig().Title,
		Players:   t.PlayerManager.GetAllPlayersState(),
	}

	event := game_session_contracts.Action{
		Type:    "full_state",
		Payload: fullStatePayload,
	}

	t.SendToPlayer(userID, event)
}

func (t *Template) EncodeMsg(action *game_session_contracts.Action, dest any) error {
	if action.Payload == nil {
		return nil
	}

	payloadBytes, err := json.Marshal(action.Payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(payloadBytes, dest)
}
