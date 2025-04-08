package referral_controller

import (
	"github.com/gofiber/fiber/v2"
	referral_service "github.com/root9464/Go_GamlerDefi/module/referral/service"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
)

type IReferralController interface {
	ReferralProcess(c *fiber.Ctx) error
	PrecheckoutReferrer(c *fiber.Ctx) error
}

type ReferralController struct {
	logger          *logger.Logger
	referralService referral_service.IReferralService
}

func NewReferralController(logger *logger.Logger, service referral_service.IReferralService) IReferralController {
	return &ReferralController{
		logger:          logger,
		referralService: service,
	}
}
