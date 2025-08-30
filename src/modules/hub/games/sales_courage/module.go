package sales_courage_module

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
	game_session_registry "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/registry"
	game_config_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/repository"
	sales_courage_dto "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/sales_courage/dto"
	deckboard_core "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/core"
	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
	"github.com/shopspring/decimal"
)

type SalesCourageSettings struct {
	HostID    string                  `mapstructure:"host_id"`
	DiceCount int                     `mapstructure:"dice_count"`
	DiceFaces int                     `mapstructure:"dice_faces"`
	Decks     []deckboard_models.Deck `mapstructure:"decks"`
}

type Game struct {
	deckboard_core.Template
	Settings SalesCourageSettings
}

func (g *Game) GetSettingsModel() any {
	return new(SalesCourageSettings)
}

func init() {
	gameProvider := &Game{}
	factory := func(repo game_config_repository.GameConfigRepository) (game_session_contracts.Game, error) {
		config, err := repo.GetByName(context.Background(), "sales_courage")

		if err != nil {
			return nil, fmt.Errorf("failed to get game config: %v", err)
		}

		var settings SalesCourageSettings
		if err := mapstructure.Decode(config.Settings, &settings); err != nil {
			return nil, fmt.Errorf("failed to decode game config: %v", err)
		}

		game := &Game{Settings: settings}

		return game, nil

	}

	game_session_registry.RegisterGame("sales_courage", factory, gameProvider)
}

func (g *Game) GetGameConfig() deckboard_models.GameConfig {
	// decks := []deckboard_models.Deck{{
	// 	ID:   "deck_0",
	// 	Name: "Deck 0",
	// 	Cards: []deckboard_models.Card{{
	// 		ID:          "card_0_0",
	// 		Title:       "Card 0_0",
	// 		Description: "Description 0_0",
	// 		Task:        "Task 0_0",
	// 		ImageURL:    "https://via.placeholder.com/150",
	// 	}},
	// }}
	//
	// return deckboard_models.GameConfig{
	// 	Title:       "Sales Courage",
	// 	Description: "Description",
	// 	Decks:       decks,
	// 	HostID:      "host_id",
	// 	DiceCount:   1,
	// 	DiceFaces:   6,
	// }

	return deckboard_models.GameConfig{
		Title:       "Sales Courage",
		Description: "Description",
		Decks:       g.Settings.Decks,
		HostID:      g.Settings.HostID,
		DiceCount:   g.Settings.DiceCount,
		DiceFaces:   g.Settings.DiceFaces,
	}
}

func (g *Game) Initialize(
	hostID string,
	sendToAll func(message any),
	sendToPlayer func(userID string, message any) error,
	broadcastToAllExcept func(excludeUserID string, message any),
) {
	g.Template.Inizialize(hostID, sendToAll, sendToPlayer, broadcastToAllExcept, g)
}

func (g *Game) HandleAction(clientID string, action game_session_contracts.Action) {
	switch action.Type {
	case "add_coins":
		g.HandleHostAction(clientID, &action, func() {
			data := new(sales_courage_dto.AddCoins)
			if err := g.Template.EncodeMsg(&action, data); err != nil {
				return
			}

			player, err := g.Template.PlayerManager.GetPlayer(data.PlayerID)
			if err != nil {
				return
			}

			g.PlayerManager.IncrementMetadataValue(player.ID, "coins", decimal.NewFromInt(int64(data.Coins)))
		})
	}

	g.Template.HandleAction(clientID, &action)
}
