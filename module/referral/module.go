package referral_module

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/root9464/Go_GamlerDefi/config"
	referral_controller "github.com/root9464/Go_GamlerDefi/module/referral/controller"
	referral_helper "github.com/root9464/Go_GamlerDefi/module/referral/helpers"
	referral_service "github.com/root9464/Go_GamlerDefi/module/referral/service"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/ton"
)

type ReferralModule struct {
	config     *config.Config
	logger     *logger.Logger
	validator  *validator.Validate
	ton_client *ton.APIClient
	ton_api    *tonapi.Client

	referralController referral_controller.IReferralController
	referralService    referral_service.IReferralService
	refferalHelper     referral_helper.IReferralHelper
}

func NewReferralModule(config *config.Config, logger *logger.Logger, validator *validator.Validate, ton_client *ton.APIClient, ton_api *tonapi.Client) *ReferralModule {
	return &ReferralModule{
		config:     config,
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
		m.referralService = referral_service.NewReferralService(m.logger, m.ton_client, m.ton_api, m.config, m.Helper())
	}
	return m.referralService
}

func (m *ReferralModule) Helper() referral_helper.IReferralHelper {
	if m.refferalHelper == nil {
		m.refferalHelper = referral_helper.NewReferralHelper(m.logger, m.config.SmartContractJettonWallet)
	}
	return m.refferalHelper
}

func (m *ReferralModule) RegisterRoutes(app fiber.Router) {
	referral := app.Group("/referral")
	referral.Post("/from-platform", m.Controller().ReferralProcessPlatform)
	referral.Get("/precheckout/:user_id", m.Controller().PrecheckoutReferrer)
}
