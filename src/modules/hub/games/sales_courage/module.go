package sales_courage_module

import (
	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
	game_session_registry "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/registry"
	sales_courage_dto "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/sales_courage/dto"
	deckboard_core "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/core"
	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
)

type SalesCourageModule struct {
}

type Game struct {
	deckboard_core.Template
}

func (g *Game) GetGameConfig() deckboard_models.GameConfig {
	decks := []deckboard_models.Deck{{
		ID:   "deck_0",
		Name: "Deck 0",
		Cards: []deckboard_models.Card{{
			ID:          "card_0_0",
			Title:       "Card 0_0",
			Description: "Description 0_0",
			Task:        "Task 0_0",
			ImageURL:    "https://via.placeholder.com/150",
		}},
	}}

	return deckboard_models.GameConfig{
		Title:       "Sales Courage",
		Description: "Description",
		Decks:       decks,
		HostID:      "host_id",
		DiceCount:   1,
		DiceFaces:   6,
	}
}

func init() {
	game_session_registry.RegisterGame("sales_courage", func() game_session_contracts.Game {
		return &Game{}
	})
}

func (g *Game) Initialize(
	sendToAll func(message any),
	sendToPlayer func(userID string, message any) error,
	broadcastToAllExcept func(excludeUserID string, message any),
) {
	g.Template.Inizialize(sendToAll, sendToPlayer, broadcastToAllExcept, g)
}

func (g *Game) HandleAction(clientID string, action game_session_contracts.Action) {
	switch action.Type {
	case "add_coins":
		data := new(sales_courage_dto.AddCoins)
		if err := g.Template.EncodeMsg(&action, data); err != nil {
			return
		}

		player, err := g.Template.PlayerManager.GetPlayer(data.PlayerID)
		if err != nil {
			return
		}

		g.PlayerManager.IncrementMetadataValue(player.ID, "coins", data.Coins)
	}

	g.Template.HandleAction(clientID, &action)
}
