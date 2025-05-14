package referral_controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	referral_repository "github.com/root9464/Go_GamlerDefi/module/referral/repository"
	referral_service "github.com/root9464/Go_GamlerDefi/module/referral/service"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
)

type IReferralController interface {
	ReferralProcessPlatform(c *fiber.Ctx) error
	PrecheckoutReferrer(c *fiber.Ctx) error
	GetDebtAuthor(c *fiber.Ctx) error
	DeletePaymentOrder(c *fiber.Ctx) error
	PayDebtAuthor(c *fiber.Ctx) error
}

type ReferralController struct {
	logger    *logger.Logger
	validator *validator.Validate

	referral_service    referral_service.IReferralService
	referral_repository referral_repository.IReferralRepository
}

func NewReferralController(
	logger *logger.Logger, validator *validator.Validate,
	referral_service referral_service.IReferralService, referral_repository referral_repository.IReferralRepository,
) IReferralController {
	return &ReferralController{
		logger:              logger,
		validator:           validator,
		referral_service:    referral_service,
		referral_repository: referral_repository,
	}
}
