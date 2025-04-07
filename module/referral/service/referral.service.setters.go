package referral_service

import (
	"context"

	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
)

func (s *ReferralService) CalculateReferralBonuses(ctx context.Context, referral referral_dto.ReferralResponse) error {
	s.logger.Infof("CalculateReferralBonuses: %+v", referral)

	return nil
}
