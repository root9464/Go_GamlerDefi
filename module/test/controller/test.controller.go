package test_controller

import (
	"github.com/gofiber/fiber/v2"
	test_service "github.com/root9464/Go_GamlerDefi/module/test/service"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
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

func (c *testController) Ping(ctx *fiber.Ctx) error {
	message := c.service.Ping()
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": message})
}
