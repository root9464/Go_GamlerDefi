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
	url      = "https://serv.gamler.atma-dev.ru/referral"
	maxLevel = 2
)

func (s *ReferralService) AccrueUserBonus(userID int, amount int) error {
	s.logger.Infof("accruing bonus for user %d: %d", userID, amount)
	_, err := utils.Patch[referral_dto.ChangeBalanceUserResponse](
		fmt.Sprintf("%s/user/%d/balance", url, userID),
		referral_dto.ChangeBalanceUserRequest{Amount: amount},
	)
	if err != nil {
		s.logger.Errorf("failed to accrue bonus for user %d: %v", userID, err)
		return err
	}
	return nil
}

func (s *ReferralService) GetReferrerChain(userID int) (*referral_dto.ReferrerResponse, error) {
	s.logger.Infof("fetching referrer chain for user %d", userID)
	resp, err := utils.Get[referral_dto.ReferrerResponse](
		fmt.Sprintf("%s/referrer/%d", url, userID),
	)
	if err != nil {
		s.logger.Errorf("failed to fetch referrer chain: %v", err)
		return nil, err
	}
	return &resp, nil
}

func (s *ReferralService) CalculateReferralBonuses(ctx context.Context, req referral_dto.ReferralProcessRequest) error {
	s.logger.Infof("starting referral bonus calculation for: %+v", req)

	bonusRates := map[int]float64{
		0: 0.20, // Уровень 1: 20%
		1: 0.02, // Уровень 2: 2%
	}

	s.logger.Infof("bonusRates: %+v", bonusRates)
	s.logger.Infof("req.PaymentType: %+v", req.PaymentType)

	switch req.PaymentType {
	case referral_dto.PaymentAuthor:
		s.logger.Infof("req.ReferrerID: %+v", req.ReferrerID)
		referrerL1, err := s.GetReferrerChain(req.ReferrerID)
		if err != nil {
			s.logger.Errorf("failed to fetch first level referrer for user %d: %v", req.ReferrerID, err)
			return errors.NewError(500, "failed to fetch referrer chain")
		}
		s.logger.Infof("referrerL1: %+v", referrerL1)

		if !lo.ContainsBy(referrerL1.ReferredUsers, func(u referral_dto.ReferredUserResponse) bool { return u.UserID == req.ReferredID }) {
			s.logger.Warnf("invalid first level referral: %+v", req.ReferredID)
			return errors.NewError(400, "invalid first level referral")
		}

		s.logger.Infof("req.ReferredID: %+v", req.ReferredID)
		s.logger.Infof("accruing bonus for levels")
		for level := 0; level <= maxLevel; level++ {
			s.logger.Infof("accruing bonus for level %d", level)
			rate, ok := bonusRates[level]
			if !ok {
				s.logger.Warnf("No bonus rate for level %d", level)
				break
			}

			s.logger.Infof("rate for level %d: %+v", level, rate)
			bonusValue := int(math.Round(float64(req.TicketCount)*rate*100) / 100)

			s.logger.Infof("bonusValue for level %d: %+v", level, bonusValue)

			if err := s.AccrueUserBonus(req.ReferrerID, bonusValue); err != nil {
				s.logger.Errorf("level %d bonus error: %v", level, err)
				return errors.NewError(500, "bonus accrual failed")
			}

			s.logger.Infof("level %d bonus accrued: %d to %d", level, bonusValue, req.ReferrerID)

			parentData, err := s.GetReferrerChain(req.ReferrerID)
			if err != nil || parentData.ReferID == 0 {
				s.logger.Warnf("stopping referral chain for user %d at level %d", req.ReferrerID, level+1)
				break
			}
			req.ReferrerID = parentData.ReferID
		}

	case referral_dto.PaymentReferred:
		s.logger.Errorf("unsupported payment type: %s", req.PaymentType)
		return errors.NewError(501, "payment type not implemented")

	default:
		s.logger.Errorf("invalid payment type: %s", req.PaymentType)
		return errors.NewError(400, "invalid payment type")
	}

	return nil
}
