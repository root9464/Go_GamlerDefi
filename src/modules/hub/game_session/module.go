package game_session_module

import (
	"github.com/gofiber/fiber/v2"
	game_session_ws "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/delivery/ws"
	game_session_hub "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/hub.go"
	_ "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/sales_courage"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type GameSessionModule struct {
	game_session_hub *game_session_hub.Hub
	game_session_ws  *game_session_ws.GameSessionHandler
	logger           *logger.Logger
}

func NewGameSessionModule(logger *logger.Logger) *GameSessionModule {
	return &GameSessionModule{
		logger: logger,
	}
}

func (m *GameSessionModule) init() {
	m.game_session_hub = game_session_hub.NewHub(m.logger)
	m.game_session_ws = game_session_ws.NewGameSessionHandler(m.game_session_hub, m.logger)
}

func (m *GameSessionModule) InitDelivery(router fiber.Router) {
	m.init()
	m.game_session_ws.RegisterRoutes(router)
}
