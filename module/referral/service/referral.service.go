package referral_service

import (
	"context"

	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/ton"
)

type IReferralService interface {
	CalculateReferralBonuses(ctx context.Context, referrer referral_dto.ReferralProcessRequest) error
}

type ReferralService struct {
	logger                       *logger.Logger
	ton_client                   *ton.APIClient
	ton_api                      *tonapi.Client
	platform_smart_contract      string
	smart_contract_jetton_wallet string
}

func NewReferralService(logger *logger.Logger, ton_client *ton.APIClient, ton_api *tonapi.Client, platform_smart_contract string, smart_contract_jetton_wallet string) IReferralService {
	return &ReferralService{
		logger:                       logger,
		ton_client:                   ton_client,
		ton_api:                      ton_api,
		platform_smart_contract:      platform_smart_contract,
		smart_contract_jetton_wallet: smart_contract_jetton_wallet,
	}
}
