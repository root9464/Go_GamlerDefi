package referral_service

import (
	"context"

	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
)

type IReferralService interface {
	CalculateReferralBonuses(ctx context.Context, referral referral_dto.ReferralResponse) error
}

type ReferralService struct {
	logger *logger.Logger
}

func NewReferralService(logger *logger.Logger) IReferralService {
	return &ReferralService{logger: logger}
}
