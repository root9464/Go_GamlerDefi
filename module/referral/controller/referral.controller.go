package referral_controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	referral_service "github.com/root9464/Go_GamlerDefi/module/referral/service"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
)

type IReferralController interface {
	// ReferralProcess(c *fiber.Ctx) error
	ReferralProcessPlatform(c *fiber.Ctx) error
	PrecheckoutReferrer(c *fiber.Ctx) error
}

type ReferralController struct {
	logger    *logger.Logger
	validator *validator.Validate

	referralService referral_service.IReferralService
}

func NewReferralController(logger *logger.Logger, validator *validator.Validate, service referral_service.IReferralService) IReferralController {
	return &ReferralController{
		logger:          logger,
		validator:       validator,
		referralService: service,
	}
}
