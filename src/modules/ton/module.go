package ton_module

import (
	"github.com/gofiber/fiber/v2"
	"github.com/root9464/Go_GamlerDefi/src/config"
	ton_controllers "github.com/root9464/Go_GamlerDefi/src/modules/ton/controllers"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type TonModule struct {
	config *config.Config
	logger *logger.Logger

	ton_controller ton_controllers.ITonController
}

func NewTonModule(config *config.Config, logger *logger.Logger) *TonModule {
	return &TonModule{config: config, logger: logger}
}

func (m *TonModule) Controller() ton_controllers.ITonController {
	if m.ton_controller == nil {
		m.ton_controller = ton_controllers.NewTonController(m.logger)
	}
	return m.ton_controller
}

func (m *TonModule) RegisterRoutes(app fiber.Router) {
	ton := app.Group("/ton")
	ton.Get("/image/:image_path", m.Controller().GetImage)
	ton.Get("/manifest", m.Controller().GetManifest)
}
