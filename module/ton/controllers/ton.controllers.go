package ton_controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
)

type ITonController interface {
	GetImage(c *fiber.Ctx) error
	GetManifest(c *fiber.Ctx) error
}

type TonController struct {
	logger *logger.Logger
}

func NewTonController(logger *logger.Logger) ITonController {
	return &TonController{logger: logger}
}
