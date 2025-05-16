package validation_module

import (
	"github.com/go-playground/validator/v10"
	"github.com/root9464/Go_GamlerDefi/config"
	validation_repository "github.com/root9464/Go_GamlerDefi/module/validation/repository"
	validation_service "github.com/root9464/Go_GamlerDefi/module/validation/service"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ValidationModule struct {
	config    *config.Config
	logger    *logger.Logger
	validator *validator.Validate
	db        *mongo.Database

	validation_service    validation_service.IValidationService
	validation_repository validation_repository.IValidationRepository
}
