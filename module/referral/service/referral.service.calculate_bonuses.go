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

type ReferralLevel struct {
	Level      int
	Rate       float64
	ReferrerID int
	Err        error
}

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

func (s *ReferralService) ReferralChainIterator(startReferrerID int, bonusRates map[int]float64, maxLevel int) iter.Seq[ReferralLevel] {
	return func(yield func(ReferralLevel) bool) {
		currentReferrerID := startReferrerID

		s.logger.Infof("starting referral chain iterator for user %d", currentReferrerID)
		for level := 0; level <= maxLevel; level++ {
			s.logger.Infof("level: %d", level)
			rate, ok := bonusRates[level]
			if !ok {
				s.logger.Warnf("no bonus rate for level %d", level)
				return
			}
			s.logger.Infof("rate: %f", rate)
			if currentReferrerID == 0 {
				s.logger.Warnf("stopping referral chain at level %d", level+1)
				return
			}

			s.logger.Infof("yielding level %d: %d", level, currentReferrerID)
			if !yield(ReferralLevel{Level: level, Rate: rate, ReferrerID: currentReferrerID, Err: nil}) {
				s.logger.Infof("stopping referral chain at level %d", level+1)
				return
			}

			s.logger.Infof("fetching referrer chain for user %d at level %d", currentReferrerID, level)
			parentData, err := s.GetReferrerChain(currentReferrerID)
			if err != nil {
				s.logger.Errorf("Failed to fetch referrer chain for user %d at level %d: %v", currentReferrerID, level, err)
				yield(ReferralLevel{Level: level, Rate: rate, ReferrerID: currentReferrerID, Err: err})
				return
			}
			s.logger.Infof("parentData: %+v", parentData)

			s.logger.Infof("referrer chain for user %d at level %d: %+v", currentReferrerID, level, parentData)
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
	s.logger.Infof("req.PaymentType: %+v | req.ReferrerID: %+v | req.ReferredID: %+v", req.PaymentType, req.ReferrerID, req.ReferredID)

	s.logger.Infof("fetching first level referrer for user %d", req.ReferrerID)
	referrerL1, err := s.GetReferrerChain(req.ReferrerID)
	if err != nil {
		s.logger.Errorf("failed to fetch first level referrer for user %d: %v", req.ReferrerID, err)
		return errors.NewError(500, "failed to fetch referrer chain")
	}
	s.logger.Infof("referrer L1: %+v", referrerL1)

	if !lo.ContainsBy(referrerL1.ReferredUsers, func(u referral_dto.ReferredUserResponse) bool { return u.UserID == req.ReferredID }) {
		s.logger.Warnf("invalid first level referral: %+v", req.ReferredID)
		return errors.NewError(400, "invalid first level referral")
	}

	switch req.PaymentType {
	case referral_dto.PaymentAuthor:
		s.logger.Infof("req.ReferredID: %+v | req.ReferrerID: %+v | req.TicketCount: %+v", req.ReferredID, req.ReferrerID, req.TicketCount)
		s.logger.Infof("accruing bonus for levels maxLevel: %d", maxLevel)
		for referralLevel := range s.ReferralChainIterator(req.ReferrerID, bonusRates, maxLevel) {
			if referralLevel.Err != nil {
				s.logger.Errorf("error in referral chain at level %d: %v", referralLevel.Level, referralLevel.Err)
				return errors.NewError(500, "error in referral chain")
			}
			level := referralLevel.Level
			rate := referralLevel.Rate
			referrerID := referralLevel.ReferrerID

			bonusValue := int(math.Round(float64(req.TicketCount)*rate*100) / 100)

			s.logger.Infof("accruing level %d bonus: %d to referrer %d", level, bonusValue, referrerID)
			if err := s.ChangeUserBalance(referrerID, bonusValue); err != nil {
				s.logger.Errorf("level %d bonus error: %v", level, err)
				return errors.NewError(500, "bonus accrual failed")
			}
		}
		return nil
	case referral_dto.PaymentReferred:
		s.logger.Infof("req.ReferredID: %+v | req.ReferrerID: %+v | req.AuthorID: %+v | req.TicketCount: %+v", req.ReferredID, req.ReferrerID, req.AuthorID, req.TicketCount)

		if req.AuthorID == 0 {
			s.logger.Warnf("author ID is required for payment type %s", req.PaymentType)
			return errors.NewError(400, "author ID is required for payment type referred")
		}

		totalBonus := 0

		s.logger.Infof("calculating total bonus for levels maxLevel: %d", maxLevel)
		for referralLevel := range s.ReferralChainIterator(req.ReferrerID, bonusRates, maxLevel) {
			if referralLevel.Err != nil {
				s.logger.Errorf("error in referral chain at level %d: %v", referralLevel.Level, referralLevel.Err)
				return errors.NewError(500, "error in referral chain")
			}
			totalBonus += int(math.Round(float64(req.TicketCount) * referralLevel.Rate))
		}

		s.logger.Infof("total bonus: %d", totalBonus)
		s.logger.Infof("debiting balance %d for author %d", totalBonus, req.AuthorID)
		if err := s.ChangeUserBalance(req.AuthorID, -totalBonus); err != nil {
			s.logger.Errorf("failed to debit balance for author %d: %v", req.AuthorID, err)
			return errors.NewError(500, "failed to debit balance for author")
		}

		s.logger.Infof("accruing bonus for levels maxLevel: %d", maxLevel)
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
