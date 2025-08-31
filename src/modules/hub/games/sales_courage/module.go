package sales_courage_module

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
	game_session_registry "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/registry"
	game_config_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/repository"
	sales_courage_dto "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/sales_courage/dto"
	deckboard_core "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/core"
	deckboard_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/submodules/deckboard/models"
	"github.com/shopspring/decimal"
)

type SalesCourageSettings struct {
	HostID    string                  `json:"host_id" mapstructure:"host_id"`
	DiceCount int                     `json:"dice_count" mapstructure:"dice_count"`
	DiceFaces int                     `json:"dice_faces" mapstructure:"dice_faces"`
	Decks     []deckboard_models.Deck `json:"decks" mapstructure:"decks"`
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

		log.Warnf("Config = %+v", config)

		var settings SalesCourageSettings

		// --- НАЧАЛО ИЗМЕНЕНИЙ ---

		// 1. Сериализуем map[string]interface{} (с BSON-типами) в стандартный JSON
		jsonBytes, err := json.Marshal(config.Settings)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal settings to json: %v", err)
		}

		// 2. Десериализуем чистый JSON в нашу целевую Go-структуру
		if err := json.Unmarshal(jsonBytes, &settings); err != nil {
			return nil, fmt.Errorf("failed to unmarshal settings from json: %v", err)
		}

		// --- КОНЕЦ ИЗМЕНЕНИЙ ---

		log.Warnf("Settings = %+v", settings)
		game := &Game{Settings: settings}
		return game, nil
	}

	game_session_registry.RegisterGame("sales_courage", factory, gameProvider)
}

func (g *Game) GetGameConfig() deckboard_models.GameConfig {
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
