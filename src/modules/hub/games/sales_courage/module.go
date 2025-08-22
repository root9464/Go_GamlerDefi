package sales_courage_module

import (
	"log"

	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
	game_session_registry "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/registry"
	deckboard_core "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/core"
	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
)

type SalesCourageModule struct {
}

type Game struct {
	deckboard_core.Template
}

func (g *Game) GetGameConfig() deckboard_models.GameConfig {
	decks := []deckboard_models.Deck{}
	for i := 0; i < 3; i++ {
		decks = append(decks, deckboard_models.Deck{
			ID:    "deck_" + string(i),
			Name:  "Deck " + string(i),
			Cards: []deckboard_models.Card{},
		})

		for j := 0; j < 3; j++ {
			decks[i].Cards = append(decks[i].Cards, deckboard_models.Card{
				ID:          "card_" + string(i) + "_" + string(j),
				Title:       "Card " + string(i) + "_" + string(j),
				Description: "Description " + string(i) + "_" + string(j),
				Task:        "Task " + string(i) + "_" + string(j),
				ImageURL:    "https://via.placeholder.com/150",
			})
		}
	}

	return deckboard_models.GameConfig{
		Title:        "Sales Courage",
		Description:  "Description",
		Decks:        []deckboard_models.Deck{},
		MaxDiceCount: 3,
	}
}

func init() {
	game_session_registry.RegisterGame("sales_courage", func() game_session_contracts.Game {
		return &Game{}
	})
}

func (g *Game) Initialize(sessID string, sendToAll func(message any), sendToPlayer func(userID string, message any) error) {
	g.Template.Inizialize(sendToAll, sendToPlayer, g)
}

func (g *Game) HandleAction(clientID string, action game_session_contracts.Action) {
	if action.Type == "make_bet" {
		log.Printf("Client %s made a bet: %v", clientID, action.Payload)
	}
	g.Template.HandleAction(clientID, &action)
}
