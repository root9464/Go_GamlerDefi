package referral_service

import (
	"context"

	"github.com/root9464/Go_GamlerDefi/config"
	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	referral_helper "github.com/root9464/Go_GamlerDefi/module/referral/helpers"
	referral_repository "github.com/root9464/Go_GamlerDefi/module/referral/repository"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/shopspring/decimal"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/ton"
)

type IReferralService interface {
	ReferralProcess(ctx context.Context, referrer referral_dto.ReferralProcessRequest) error
	PayPaymentOrder(ctx context.Context, paymentOrderID string, walletAddress string) (string, error)
	PayAllPaymentOrders(ctx context.Context, authorID int, walletAddress string) (string, error)
	AssessInvitationAbility(ctx context.Context, authorID int) (bool, error)
	CalculateAuthorDebt(ctx context.Context, authorID int) (decimal.Decimal, error)
}

type ReferralService struct {
	logger     *logger.Logger
	ton_client *ton.APIClient
	ton_api    *tonapi.Client
	config     *config.Config

	referral_helper     referral_helper.IReferralHelper
	referral_repository referral_repository.IReferralRepository
}

func NewReferralService(
	logger *logger.Logger,
	ton_client *ton.APIClient,
	ton_api *tonapi.Client,
	config *config.Config,

	referral_helper referral_helper.IReferralHelper,
	referral_repository referral_repository.IReferralRepository,
) IReferralService {
	return &ReferralService{
		logger:              logger,
		ton_client:          ton_client,
		ton_api:             ton_api,
		config:              config,
		referral_helper:     referral_helper,
		referral_repository: referral_repository,
	}
}
