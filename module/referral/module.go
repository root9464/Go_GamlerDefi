package referral_module

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	referral_controller "github.com/root9464/Go_GamlerDefi/module/referral/controller"
	referral_service "github.com/root9464/Go_GamlerDefi/module/referral/service"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/ton"
)

type ReferralModule struct {
	logger     *logger.Logger
	validator  *validator.Validate
	ton_client *ton.APIClient
	ton_api    *tonapi.Client

	referralController referral_controller.IReferralController
	referralService    referral_service.IReferralService
}

func NewReferralModule(logger *logger.Logger, validator *validator.Validate, ton_client *ton.APIClient, ton_api *tonapi.Client) *ReferralModule {
	return &ReferralModule{
		logger:     logger,
		validator:  validator,
		ton_client: ton_client,
		ton_api:    ton_api,
	}
}

func (m *ReferralModule) Controller() referral_controller.IReferralController {
	if m.referralController == nil {
		m.referralController = referral_controller.NewReferralController(m.logger, m.validator, m.Service())
	}
	return m.referralController
}

func (m *ReferralModule) Service() referral_service.IReferralService {
	if m.referralService == nil {
		m.referralService = referral_service.NewReferralService(m.logger, m.ton_client)
	}
	return m.referralService
}

func (m *ReferralModule) RegisterRoutes(app fiber.Router) {
	referral := app.Group("/referral")
	referral.Post("/from-platform", m.Controller().ReferralProcessPlatform)
	referral.Get("/precheckout/:user_id", m.Controller().PrecheckoutReferrer)
}
