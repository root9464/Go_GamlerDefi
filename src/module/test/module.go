package test_module

import (
	"github.com/gofiber/fiber/v2"
	test_controller "github.com/root9464/Go_GamlerDefi/src/module/test/controller"
	test_service "github.com/root9464/Go_GamlerDefi/src/module/test/service"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type TestModule struct {
	logger *logger.Logger

	controller test_controller.ITestController
	service    test_service.ITestService
}

func NewTestModule(logger *logger.Logger) *TestModule {
	return &TestModule{
		logger: logger,
	}
}

func (m *TestModule) Controller() test_controller.ITestController {
	if m.controller == nil {
		m.controller = test_controller.NewTestController(m.logger, m.Service())
	}
	return m.controller
}

func (m *TestModule) Service() test_service.ITestService {
	if m.service == nil {
		m.service = test_service.NewTestService(m.logger)
	}
	return m.service
}

func (m *TestModule) RegisterRoutes(app fiber.Router) {
	app.Get("/ping", m.Controller().Ping)
}
