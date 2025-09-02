package game_session_module

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/root9464/Go_GamlerDefi/src/config"
	conference_ws "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/handler"
	game_session_http "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/delivery/http"
	game_session_ws "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/delivery/ws"
	game_session_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/infrastructure/repository/game_session"
	trash_repo "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/infrastructure/trash/repository"
	game_session_hub "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/hub.go"
	game_config_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/repository"
	_ "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/sales_courage"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

type GameSessionModule struct {
	game_session_repository *game_session_repository.GameSessionRepository
	trash_repository        *trash_repo.TrashRepository
	game_config_repository  *game_config_repository.GameConfigRepository

	game_session_hub *game_session_hub.Hub

	game_session_ws    *game_session_ws.GameSessionHandler
	game_seession_http *game_session_http.Handler
	conference_ws      *conference_ws.WSHandler

	logger   *logger.Logger
	db       *mongo.Database
	postgres *gorm.DB
	config   *config.Config
}

func NewGameSessionModule(
	logger *logger.Logger,
	db *mongo.Database,
	postgres *gorm.DB,
	gameConfigRepository *game_config_repository.GameConfigRepository,
	conference_ws *conference_ws.WSHandler,
	config *config.Config,
) *GameSessionModule {
	log.Warnf("Conf ws = %v", conference_ws)

	return &GameSessionModule{
		logger:   logger,
		db:       db,
		postgres: postgres,
		config:   config,

		game_config_repository: gameConfigRepository,

		conference_ws: conference_ws,
	}
}

func (m *GameSessionModule) init() {
	m.game_session_repository = game_session_repository.NewGameSessionRepository(m.logger, m.db)
	m.trash_repository = trash_repo.NewTrashTicketRepository(m.logger, m.postgres)

	m.game_session_hub = game_session_hub.NewHub(m.logger, m.game_session_repository, m.trash_repository, m.game_config_repository)

	m.game_session_ws = game_session_ws.NewGameSessionHandler(m.game_session_hub, m.logger, m.conference_ws, m.config)
	m.game_seession_http = game_session_http.NewGameSessionHandler(m.game_session_hub, m.logger)
}

func (m *GameSessionModule) InitDelivery(router fiber.Router) {
	m.init()
	m.game_session_ws.RegisterRoutes(router)
	m.game_seession_http.RegisterRoutes(router)
}
