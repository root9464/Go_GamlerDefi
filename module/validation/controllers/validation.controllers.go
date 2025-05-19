package validation_controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	validation_service "github.com/root9464/Go_GamlerDefi/module/validation/service"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
)

type IValidationController interface {
	ValidatorTransaction(c *fiber.Ctx) error
}

type ValidationController struct {
	logger    *logger.Logger
	validator *validator.Validate

	validation_service validation_service.IValidationService
}

func NewValidationController(logger *logger.Logger, validator *validator.Validate, validation_service validation_service.IValidationService) IValidationController {
	return &ValidationController{logger: logger, validator: validator, validation_service: validation_service}
}
