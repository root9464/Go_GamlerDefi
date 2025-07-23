package test_controller

import (
	"github.com/gofiber/fiber/v2"
	test_service "github.com/root9464/Go_GamlerDefi/src/modules/test/service"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type ITestController interface {
	Ping(ctx *fiber.Ctx) error
}

type testController struct {
	logger  *logger.Logger
	service test_service.ITestService
}

func NewTestController(logger *logger.Logger, service test_service.ITestService) *testController {
	return &testController{logger: logger, service: service}
}

// @Summary Ping
// @Description Ping the server
// @Tags test
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/ping [get]
func (c *testController) Ping(ctx *fiber.Ctx) error {
	message := c.service.Ping()
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": message})
}
