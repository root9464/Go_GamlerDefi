package referral_service

import (
	"context"
	"fmt"
	"math"

	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"github.com/root9464/Go_GamlerDefi/packages/utils"
	"github.com/samber/lo"
)

const (
	url = "https://serv.gamler.atma-dev.ru/referral"
)

func (s *ReferralService) CalculateReferralBonuses(ctx context.Context, req referral_dto.ReferralProcessRequest) error {
	s.logger.Infof("Calculating bonuses for: %+v", req)

	response, err := utils.Get[referral_dto.ReferrerResponse](fmt.Sprintf("%s/referrer/%d", url, req.ReferrerID))
	if err != nil {
		s.logger.Errorf("Fetch error: %v", err)
		return errors.NewError(404, err.Error())
	}
	if len(response.ReferredUsers) == 0 {
		return fmt.Errorf("no referred users found")
	}
	if req.PaymentType != referral_dto.PaymentAuthor {
		return errors.NewError(400, "invalid payment type")
	}

	user, exists := lo.Find(response.ReferredUsers, func(u referral_dto.ReferredUserResponse) bool {
		return u.UserID == req.ReferredID
	})
	if !exists {
		return fmt.Errorf("referred user %d not found", req.ReferredID)
	}

	index := lo.IndexOf(response.ReferredUsers, user)
	if index < 0 || index > 1 {
		return nil
	}

	rates := []float64{0.20, 0.02}
	bonus := math.Round(float64(req.TicketCount)*rates[index]*100) / 100

	s.logger.Infof("accrual of bonuses under the referral program: %.2f", bonus)

	return nil
}
