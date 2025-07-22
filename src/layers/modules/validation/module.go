package validation_module

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/root9464/Go_GamlerDefi/src/config"
	validation_controllers "github.com/root9464/Go_GamlerDefi/src/layers/modules/validation/controllers"
	validation_repository "github.com/root9464/Go_GamlerDefi/src/layers/modules/validation/repository"
	validation_service "github.com/root9464/Go_GamlerDefi/src/layers/modules/validation/service"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"github.com/tonkeeper/tonapi-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ValidationModule struct {
	config    *config.Config
	logger    *logger.Logger
	validator *validator.Validate
	db        *mongo.Database
	ton_api   *tonapi.Client

	validation_service    validation_service.IValidationService
	validation_repository validation_repository.IValidationRepository
	validation_controller validation_controllers.IValidationController
}

func NewValidationModule(config *config.Config, logger *logger.Logger, validator *validator.Validate, db *mongo.Database, ton_api *tonapi.Client) *ValidationModule {
	return &ValidationModule{config: config, logger: logger, validator: validator, db: db, ton_api: ton_api}
}

func (m *ValidationModule) Controller() validation_controllers.IValidationController {
	if m.validation_controller == nil {
		m.validation_controller = validation_controllers.NewValidationController(m.logger, m.validator, m.Service())
	}
	return m.validation_controller
}

func (m *ValidationModule) Service() validation_service.IValidationService {
	if m.validation_service == nil {
		m.validation_service = validation_service.NewValidationService(m.logger, m.ton_api, m.Repository())
	}
	return m.validation_service
}

func (m *ValidationModule) Repository() validation_repository.IValidationRepository {
	if m.validation_repository == nil {
		m.validation_repository = validation_repository.NewValidationRepository(m.logger, m.db)
	}
	return m.validation_repository
}

func (m *ValidationModule) RegisterRoutes(app fiber.Router) {
	validation := app.Group("/validation")
	validation.Post("/validate", m.Controller().ValidatorTransaction)
}
