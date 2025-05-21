package referral_module

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/root9464/Go_GamlerDefi/config"
	referral_controller "github.com/root9464/Go_GamlerDefi/module/referral/controller"
	referral_helper "github.com/root9464/Go_GamlerDefi/module/referral/helpers"
	referral_repository "github.com/root9464/Go_GamlerDefi/module/referral/repository"
	referral_service "github.com/root9464/Go_GamlerDefi/module/referral/service"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/ton"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ReferralModule struct {
	config    *config.Config
	logger    *logger.Logger
	validator *validator.Validate
	db        *mongo.Database

	ton_client *ton.APIClient
	ton_api    *tonapi.Client

	referral_controller referral_controller.IReferralController
	referral_service    referral_service.IReferralService
	refferal_helper     referral_helper.IReferralHelper
	referral_repository referral_repository.IReferralRepository
}

func NewReferralModule(
	config *config.Config, logger *logger.Logger, validator *validator.Validate, db *mongo.Database,
	ton_client *ton.APIClient, ton_api *tonapi.Client,
) *ReferralModule {
	return &ReferralModule{
		config:     config,
		logger:     logger,
		validator:  validator,
		db:         db,
		ton_client: ton_client,
		ton_api:    ton_api,
	}
}

func (m *ReferralModule) Controller() referral_controller.IReferralController {
	if m.referral_controller == nil {
		m.referral_controller = referral_controller.NewReferralController(m.logger, m.validator, m.Service(), m.Repository())
	}
	return m.referral_controller
}

func (m *ReferralModule) Service() referral_service.IReferralService {
	if m.referral_service == nil {
		m.referral_service = referral_service.NewReferralService(m.logger, m.ton_client, m.ton_api, m.config, m.Helper(), m.Repository())
	}
	return m.referral_service
}

func (m *ReferralModule) Helper() referral_helper.IReferralHelper {
	if m.refferal_helper == nil {
		m.refferal_helper = referral_helper.NewReferralHelper(m.logger, m.config.PlatformSmartContract)
	}
	return m.refferal_helper
}

func (m *ReferralModule) Repository() referral_repository.IReferralRepository {
	if m.referral_repository == nil {
		m.referral_repository = referral_repository.NewReferralRepository(m.logger, m.db)
	}
	return m.referral_repository
}

func (m *ReferralModule) RegisterRoutes(app fiber.Router) {
	referral := app.Group("/referral")
	referral.Post("/from-platform", m.Controller().ReferralProcessPlatform)
	referral.Get("/precheckout/:user_id", m.Controller().PrecheckoutReferrer)

	referral.Get("/:author_id/payment-orders", m.Controller().GetDebtAuthor)
	referral.Delete("/payment-orders/:order_id", m.Controller().DeletePaymentOrder)
	referral.Delete("/payment-orders/all", m.Controller().DeleteAllPaymentOrders)
	referral.Get("/payment-orders/pay", m.Controller().PayDebtAuthor)             // /payment-orders/pay?order_id=<id>
	referral.Get("/payment-orders/all", m.Controller().PayAllDebtAuthor)          // /payment-orders/all?author_id=<id>
	referral.Get("/validate-invite", m.Controller().ValidateInvitationConditions) // /validate-invite?author_id=<id>
}
