package referral_service

import (
	"context"

	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/xssnick/tonutils-go/ton"
)

type IReferralService interface {
	CalculateReferralBonuses(ctx context.Context, referrer referral_dto.ReferralProcessRequest) error
}

type ReferralService struct {
	logger     *logger.Logger
	ton_client *ton.APIClient
}

func NewReferralService(logger *logger.Logger, ton_client *ton.APIClient) IReferralService {
	return &ReferralService{logger: logger, ton_client: ton_client}
}
