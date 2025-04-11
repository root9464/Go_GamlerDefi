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
		return u.UserID == req.ReferredID
	})
	if !exists {
		return fmt.Errorf("referred user %d not found", req.ReferredID)
	}

	if req.PaymentType == referral_dto.PaymentAuthor {
		bonus20 := math.Round(float64(req.TicketCount) * 0.20)
		_, err = AccruePlatformBonus(req.ReferrerID, referral_dto.ChangeBalanceUserRequest{Amount: int(bonus20)})
		if err != nil {
			s.logger.Errorf("Failed to accrue bonus to referrer %d: %v", req.ReferrerID, err)
			return err
		}
		s.logger.Infof("Accrued %d Gamler to referrer %d", int(bonus20), req.ReferrerID)

		if response.ReferID > 0 {
			bonus2 := math.Round(float64(req.TicketCount) * 0.02)
			_, err = AccruePlatformBonus(response.ReferID, referral_dto.ChangeBalanceUserRequest{Amount: int(bonus2)})
			if err != nil {
				s.logger.Errorf("Failed to accrue bonus to upper referrer %d: %v", response.ReferID, err)
				return err
			}
			s.logger.Infof("Accrued %d Gamler to upper referrer %d", int(bonus2), response.ReferID)
		}
	} else if req.PaymentType == referral_dto.PaymentReferred {
		s.logger.Warnf("Payment type %s not implemented", req.PaymentType)
		return errors.NewError(500, "Payment type not implemented")
	}

	return nil
}
