package referral_service

import (
	"context"
	"fmt"
	"iter"
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

func (s *ReferralService) CheckBalanceForReward(userID int) (int, error) {
	s.logger.Infof("checking balance for reward for user %d", userID)
	resp, err := utils.Get[referral_dto.BalanceResponse](
		fmt.Sprintf("%s/user/%d/balance", url, userID),
	)
	if err != nil {
		s.logger.Errorf("failed to check balance for reward: %v", err)
		return 0, err
	}
	return resp.Balance, nil
}

func (s *ReferralService) ChangeUserBalance(userID int, amount int) error {
	s.logger.Infof("debiting balance for user %d: %d", userID, amount)
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

type ReferralLevel struct {
	Level      int
	Rate       float64
	ReferrerID int
	Err        error
}

func (s *ReferralService) ReferralChainIterator(startReferrerID int, bonusRates map[int]float64, maxLevel int) iter.Seq[ReferralLevel] {
	return func(yield func(ReferralLevel) bool) {
		currentReferrerID := startReferrerID
		for level := 0; level <= maxLevel; level++ {
			rate, ok := bonusRates[level]
			if !ok {
				s.logger.Warnf("No bonus rate for level %d", level)
				return
			}
			if currentReferrerID == 0 {
				s.logger.Warnf("Stopping referral chain at level %d", level+1)
				return
			}
			if !yield(ReferralLevel{Level: level, Rate: rate, ReferrerID: currentReferrerID, Err: nil}) {
				return
			}
			parentData, err := s.GetReferrerChain(currentReferrerID)
			if err != nil {
				s.logger.Errorf("Failed to fetch referrer chain for user %d at level %d: %v", currentReferrerID, level, err)
				yield(ReferralLevel{Level: level, Rate: rate, ReferrerID: currentReferrerID, Err: err})
				return
			}
			if parentData.ReferID == 0 {
				s.logger.Warnf("Stopping referral chain at level %d", level+1)
				return
			}
			currentReferrerID = parentData.ReferID
		}
	}
}

func (s *ReferralService) CalculateReferralBonuses(ctx context.Context, req referral_dto.ReferralProcessRequest) error {
	s.logger.Infof("starting referral bonus calculation for: %+v", req)

	bonusRates := map[int]float64{
		0: 0.20, // Уровень 1: 20%
		1: 0.02, // Уровень 2: 2%
	}

	s.logger.Infof("bonusRates: %+v", bonusRates)
	s.logger.Infof("req.PaymentType: %+v", req.PaymentType)
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

	switch req.PaymentType {
	case referral_dto.PaymentAuthor:
		s.logger.Infof("req.ReferredID: %+v", req.ReferredID)
		s.logger.Infof("accruing bonus for levels")
		for referralLevel := range s.ReferralChainIterator(req.ReferrerID, bonusRates, maxLevel) {
			if referralLevel.Err != nil {
				s.logger.Errorf("error in referral chain at level %d: %v", referralLevel.Level, referralLevel.Err)
				return errors.NewError(500, "error in referral chain")
			}
			level := referralLevel.Level
			rate := referralLevel.Rate
			referrerID := referralLevel.ReferrerID

			bonusValue := int(math.Round(float64(req.TicketCount)*rate*100) / 100)

			s.logger.Infof("Accruing level %d bonus: %d to referrer %d", level, bonusValue, referrerID)
			if err := s.ChangeUserBalance(referrerID, bonusValue); err != nil {
				s.logger.Errorf("level %d bonus error: %v", level, err)
				return errors.NewError(500, "bonus accrual failed")
			}
		}
		return nil
	case referral_dto.PaymentReferred:
		s.logger.Infof("req.AuthorID: %+v", req.AuthorID)
		if req.AuthorID == 0 {
			s.logger.Warnf("author ID is required for payment type %s", req.PaymentType)
			return errors.NewError(400, "author ID is required for payment type referred")
		}

		s.logger.Infof("req.ReferredID: %+v", req.ReferredID)
		balanceAuthor, err := s.CheckBalanceForReward(req.AuthorID)
		if err != nil {
			s.logger.Errorf("failed to check balance for reward: %v", err)
			return errors.NewError(500, "failed to check balance for reward")
		}
		s.logger.Infof("balance for author %d: %+v", req.AuthorID, balanceAuthor)

		totalBonus := 0

		for referralLevel := range s.ReferralChainIterator(req.ReferrerID, bonusRates, maxLevel) {
			if referralLevel.Err != nil {
				s.logger.Errorf("error in referral chain at level %d: %v", referralLevel.Level, referralLevel.Err)
				return errors.NewError(500, "error in referral chain")
			}
			totalBonus += int(math.Round(float64(req.TicketCount) * referralLevel.Rate))
		}

		s.logger.Infof("total bonus: %d", totalBonus)
		if balanceAuthor < totalBonus {
			s.logger.Errorf("balance for author %d is less than total bonus %d", req.AuthorID, totalBonus)
			return errors.NewError(400, "balance for author is less than total bonus")
		}

		s.logger.Infof("debiting balance %d for author %d", totalBonus, req.AuthorID)
		if err := s.ChangeUserBalance(req.AuthorID, -totalBonus); err != nil {
			s.logger.Errorf("failed to debit balance for author %d: %v", req.AuthorID, err)
			return errors.NewError(500, "failed to debit balance for author")
		}

		s.logger.Infof("balance for author %d after debit: %d", req.AuthorID, balanceAuthor-totalBonus)
		s.logger.Infof("accruing bonus for levels")
		for referralLevel := range s.ReferralChainIterator(req.ReferrerID, bonusRates, maxLevel) {
			if referralLevel.Err != nil {
				s.logger.Errorf("error in referral chain at level %d: %v", referralLevel.Level, referralLevel.Err)
				return errors.NewError(500, "error in referral chain")
			}
			level := referralLevel.Level
			rate := referralLevel.Rate
			referrerID := referralLevel.ReferrerID

			bonusValue := int(math.Round(float64(req.TicketCount)*rate*100) / 100)

			s.logger.Infof("Accruing level %d bonus: %d to referrer %d", level, bonusValue, referrerID)
			if err := s.ChangeUserBalance(referrerID, bonusValue); err != nil {
				s.logger.Errorf("level %d bonus error: %v", level, err)
				return errors.NewError(500, "bonus accrual failed")
			}
		}
		return nil
	default:
		s.logger.Errorf("invalid payment type: %s", req.PaymentType)
		return errors.NewError(400, "invalid payment type")
	}
}
