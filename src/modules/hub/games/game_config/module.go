package game_config_module

import (
	"github.com/gofiber/fiber/v2"
	game_config_handler "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/handler"
	game_config_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/repository"
	game_config_service "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/service"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	file_module "github.com/root9464/Go_GamlerDefi/src/submodules/file"
	file_service "github.com/root9464/Go_GamlerDefi/src/submodules/file/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GameConfigModule struct {
	Game_config_repository *game_config_repository.GameConfigRepository
	game_config_service    *game_config_service.GameConfigService
	game_config_handler    *game_config_handler.GameConfigHandler

	logger      *logger.Logger
	db          *mongo.Database
	fileService *file_service.FileService
}

func NewGameConfigModule(logger *logger.Logger, db *mongo.Database) *GameConfigModule {
	fileModule := file_module.NewFileModule(logger, "../assets/")
	fileModule.Init()
	gm := GameConfigModule{
		logger:      logger,
		db:          db,
		fileService: fileModule.File_service,
	}

	gm.init()
	return &gm
}

func (m *GameConfigModule) init() {
	m.Game_config_repository = game_config_repository.NewGameConfigRepository(m.logger, m.db)
	m.game_config_service = game_config_service.NewGameConfigService(m.logger, m.Game_config_repository)
	m.game_config_handler = game_config_handler.NewGameConfigHandler(m.logger, m.game_config_service, m.fileService)
}

func (m *GameConfigModule) InitDelivery(router fiber.Router) {
	m.game_config_handler.RegisterRoutes(router)
}
