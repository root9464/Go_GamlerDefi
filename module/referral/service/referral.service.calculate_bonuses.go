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

func AccruePlatformBonus(userID int, body referral_dto.ChangeBalanceUserRequest) (referral_dto.ChangeBalanceUserResponse, error) {
	response, err := utils.Patch[referral_dto.ChangeBalanceUserResponse](
		fmt.Sprintf("%s/user/%d/balance", url, userID),
		body,
	)
	if err != nil {
		return referral_dto.ChangeBalanceUserResponse{}, err
	}
	return response, nil
}

func (s *ReferralService) CalculateReferralBonuses(ctx context.Context, req referral_dto.ReferralProcessRequest) error {
	s.logger.Infof("Calculating bonuses for: %+v", req)

	response, err := utils.Get[referral_dto.ReferrerResponse](fmt.Sprintf("%s/referrer/%d", url, req.ReferrerID))
	if err != nil {
		s.logger.Errorf("Fetch error: %v", err)
		return errors.NewError(404, err.Error())
	}

	_, exists := lo.Find(response.ReferredUsers, func(u referral_dto.ReferredUserResponse) bool {
		s.logger.Infof("Checking if user %d exists in referred users", req.ReferredID)
		return u.UserID == req.ReferredID
	})

	if !exists {
		s.logger.Errorf("referred user %d not found", req.ReferredID)
		return errors.NewError(404, fmt.Sprintf("referred user %d not found", req.ReferredID))
	}

	bonusRates := []float64{0.20, 0.02}
	if req.PaymentType == referral_dto.PaymentAuthor {
		for level, rate := range bonusRates {
			if req.ReferrerID == 0 {
				s.logger.Warnf("Referrer ID is 0 at level %d", level+1)
				continue
			}

			bonus := math.Round(float64(req.TicketCount) * rate)
			_, err := AccruePlatformBonus(req.ReferrerID, referral_dto.ChangeBalanceUserRequest{Amount: int(bonus)})
			if err != nil {
				s.logger.Errorf("Failed to accrue bonus to referrer %d at level %d: %v", req.ReferrerID, level+1, err)
				return errors.NewError(500, err.Error())
			}
			s.logger.Infof("Accrued %d Gamler to referrer %d at level %d", int(bonus), req.ReferrerID, level+1)

			referrerResponse, err := utils.Get[referral_dto.ReferrerResponse](fmt.Sprintf("%s/referrer/%d", url, req.ReferrerID))
			if err != nil {
				s.logger.Errorf("Failed to fetch referrer info for %d: %v", req.ReferrerID, err)
				return errors.NewError(500, err.Error())
			}

			if referrerResponse.ReferID == nil {
				s.logger.Infof("No parent referrer for user %d at level %d, stopping bonus calculation", req.ReferrerID, level+1)
				break
			}
			req.ReferrerID = *referrerResponse.ReferID
		}
	} else if req.PaymentType == referral_dto.PaymentReferred {
		s.logger.Warnf("Payment type %s not implemented", req.PaymentType)
		return errors.NewError(500, "Payment type not implemented")
	}

	return nil
}
