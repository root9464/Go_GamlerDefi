package game_session_module

import (
	"github.com/gofiber/fiber/v2"
	game_session_http "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/delivery/http"
	game_session_ws "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/delivery/ws"
	game_session_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/infrastructure/repository/game_session"
	game_session_hub "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/hub.go"
	_ "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/sales_courage"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GameSessionModule struct {
	game_session_hub        *game_session_hub.Hub
	game_session_ws         *game_session_ws.GameSessionHandler
	game_session_repository *game_session_repository.GameSessionRepository
	game_seession_http      *game_session_http.Handler
	logger                  *logger.Logger
	db                      *mongo.Database
}

func NewGameSessionModule(logger *logger.Logger, db *mongo.Database) *GameSessionModule {
	return &GameSessionModule{
		logger: logger,
		db:     db,
	}
}

func (m *GameSessionModule) init() {
	m.game_session_repository = game_session_repository.NewGameSessionRepository(m.logger, m.db)
	m.game_session_hub = game_session_hub.NewHub(m.logger, m.game_session_repository)
	m.game_session_ws = game_session_ws.NewGameSessionHandler(m.game_session_hub, m.logger)
	m.game_seession_http = game_session_http.NewGameSessionHandler(m.game_session_hub, m.logger)
}

func (m *GameSessionModule) InitDelivery(router fiber.Router) {
	m.init()
	m.game_session_ws.RegisterRoutes(router)
	m.game_seession_http.RegisterRoutes(router)
}
