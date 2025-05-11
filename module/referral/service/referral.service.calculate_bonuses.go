package referral_service

import (
	"context"
	"encoding/base64"
	"fmt"
	"iter"
	"math"
	"strconv"
	"time"

	referral_adapters "github.com/root9464/Go_GamlerDefi/module/referral/adapters"
	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	referral_helper "github.com/root9464/Go_GamlerDefi/module/referral/helpers"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"github.com/root9464/Go_GamlerDefi/packages/utils"
	"github.com/samber/lo"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

const (
	url      = "https://serv.gamler.atma-dev.ru/referral"
	maxLevel = 2
)

func (s *ReferralService) getReferrerChain(userID int) (*referral_dto.ReferrerResponse, error) {
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

func (s *ReferralService) getAuthorData(authorID int) (*referral_dto.ReferrerResponse, error) {
	s.logger.Infof("fetching author data for user %d", authorID)
	resp, err := utils.Get[referral_dto.ReferrerResponse](
		fmt.Sprintf("%s/referrer/%d", url, authorID),
	)
	if err != nil {
		s.logger.Errorf("failed to fetch author data for user %d: %v", authorID, err)
		return nil, err
	}
	return &resp, nil
}

type ReferralLevel struct {
	Level         int
	Rate          float64
	ReferrerID    int
	WalletAddress string
	Err           error
}

func (s *ReferralService) referralChainIterator(req referral_dto.ReferralProcessRequest, bonusRates map[int]float64, maxLevel int) iter.Seq[ReferralLevel] {
	return func(yield func(ReferralLevel) bool) {
		currentReferrerID := req.ReferrerID
		referredID := req.ReferralID

		s.logger.Infof("fetching first level referrer for user %d", currentReferrerID)
		referrerL1, err := s.getReferrerChain(currentReferrerID)
		if err != nil {
			s.logger.Errorf("failed to fetch first level referrer for user %d: %v", currentReferrerID, err)
			yield(ReferralLevel{Level: 0, Rate: 0, ReferrerID: currentReferrerID, WalletAddress: "", Err: errors.NewError(500, "failed to fetch referrer chain")})
			return
		}

		if !lo.ContainsBy(referrerL1.ReferredUsers, func(u referral_dto.ReferredUserResponse) bool { return u.UserID == referredID }) {
			s.logger.Warnf("invalid first level referral: %+v", referredID)
			yield(ReferralLevel{Level: 0, Rate: 0, ReferrerID: currentReferrerID, WalletAddress: "", Err: errors.NewError(400, "invalid first level referral")})
			return
		}

		for level := 0; level <= maxLevel; level++ {
			rate, ok := bonusRates[level]
			if !ok {
				s.logger.Warnf("no bonus rate for level %d", level)
				return
			}

			if currentReferrerID == 0 {
				s.logger.Warnf("stopping referral chain at level %d", level+1)
				return
			}

			referrerData, err := s.getReferrerChain(currentReferrerID)
			if err != nil {
				s.logger.Errorf("failed to fetch referrer data for user %d at level %d: %v", currentReferrerID, level, err)
				yield(ReferralLevel{Level: level, Rate: rate, ReferrerID: currentReferrerID, WalletAddress: "", Err: err})
				return
			}

			if !yield(ReferralLevel{Level: level, Rate: rate, ReferrerID: currentReferrerID, WalletAddress: referrerData.WalletAddress, Err: nil}) {
				s.logger.Infof("stopping referral chain at level %d", level+1)
				return
			}

			parentData, err := s.getReferrerChain(currentReferrerID)
			if err != nil {
				s.logger.Errorf("failed to fetch referrer chain for user %d at level %d: %v", currentReferrerID, level, err)
				yield(ReferralLevel{Level: level, Rate: rate, ReferrerID: currentReferrerID, WalletAddress: "", Err: err})
				return
			}

			if parentData.ReferrerID == 0 {
				s.logger.Warnf("stopping referral chain at level %d", level+1)
				return
			}

			currentReferrerID = parentData.ReferrerID
		}
	}
}

type ReferralBonusResult struct {
	TotalBonusValue   float64
	AccrualDictionary []referral_helper.JettonEntry
	Levels            []referral_dto.LevelRequest
}

func (s *ReferralService) calculateReferralBonuses(req referral_dto.ReferralProcessRequest, bonusRates map[int]float64, maxLevel int) (ReferralBonusResult, error) {
	s.logger.Infof("calculating bonus for levels maxLevel: %d", maxLevel)
	totalBonusValue := 0.0
	accrualDictionary := []referral_helper.JettonEntry{}
	levels := []referral_dto.LevelRequest{}

	for referralLevel := range s.referralChainIterator(req, bonusRates, maxLevel) {
		if referralLevel.Err != nil {
			s.logger.Errorf("error in referral chain at level %d: %v", referralLevel.Level, referralLevel.Err)
			return ReferralBonusResult{
				TotalBonusValue:   0,
				AccrualDictionary: []referral_helper.JettonEntry{},
				Levels:            []referral_dto.LevelRequest{},
			}, referralLevel.Err
		}

		rate := referralLevel.Rate
		bonusAmount := (math.Round(float64(req.TicketCount)*rate*100) / 100)
		totalBonusValue += bonusAmount
		if referralLevel.WalletAddress != "" {
			accrualDictionary = append(accrualDictionary, referral_helper.JettonEntry{
				Address: address.MustParseAddr(referralLevel.WalletAddress),
				Amount:  uint64(bonusAmount * math.Pow10(9)),
			})

			levels = append(levels, referral_dto.LevelRequest{
				LevelNumber: referralLevel.Level,
				Rate:        referralLevel.Rate,
				Amount:      bonusAmount,
			})
		}

		s.logger.Infof("level %d: referrer %d %s gets %.2f bonus", referralLevel.Level, referralLevel.ReferrerID, referralLevel.WalletAddress, bonusAmount)
	}
	s.logger.Infof("total bonus value: %f", totalBonusValue)
	s.logger.Infof("dictionary with the values of referral bonus accruals: %+v", accrualDictionary)
	return ReferralBonusResult{
		TotalBonusValue:   totalBonusValue,
		AccrualDictionary: accrualDictionary,
		Levels:            levels,
	}, nil
}

func (s *ReferralService) precheckoutBalance(targetAddress string) (float64, error) {
	s.logger.Infof("checking the balance of a author wallet for awarding bonuses")
	contractBalance, err := s.ton_api.GetAccountJettonsBalances(context.Background(), tonapi.GetAccountJettonsBalancesParams{
		AccountID: targetAddress,
	})

	if err != nil {
		s.logger.Errorf("failed to fetch account jettons balances: %v", err)
		return 0, errors.NewError(500, "failed to fetch account jettons balances")
	}

	s.logger.Infof("find target jetton address %s in balances author wallet %s", s.config.TargetJettonMaster, targetAddress)
	foundJetton, found := lo.Find(contractBalance.Balances, func(b tonapi.JettonBalance) bool {
		rawAddr, parseErr := address.ParseRawAddr(b.Jetton.Address)
		if parseErr != nil {
			s.logger.Errorf("failed to parse wallet address: %v", parseErr)
			return false
		}
		userFriendlyAddr := rawAddr.Bounce(true).Testnet(true).String()
		s.logger.Infof("user friendly address: %s", userFriendlyAddr)
		return userFriendlyAddr == s.config.TargetJettonMaster
	})

	if !found {
		s.logger.Errorf("target jetton address %s not found in balances author wallet %s", s.config.TargetJettonMaster, targetAddress)
		return 0, errors.NewError(404, "target jetton address not found")
	}
	s.logger.Infof("target jetton address %s found in balances author wallet %s %s", s.config.TargetJettonMaster, targetAddress, foundJetton.Balance)
	s.logger.Infof("convert jetton balance to float64: %s", foundJetton.Balance)
	jettonBalance, err := strconv.ParseFloat(foundJetton.Balance, 64)
	if err != nil {
		s.logger.Errorf("failed to convert jetton balance to float64: %v", err)
		return 0, errors.NewError(500, "failed to convert jetton balance to float64")
	}

	return jettonBalance, nil
}

func (s *ReferralService) ReferralProcess(ctx context.Context, req referral_dto.ReferralProcessRequest) (string, error) {
	s.logger.Infof("starting referral bonus calculation for: %+v", req)

	bonusRates := map[int]float64{
		0: 0.20, // Уровень 1: 20%
		1: 0.02, // Уровень 2: 2%
	}

	s.logger.Infof("bonusRates: %+v", bonusRates)
	s.logger.Infof("req.PaymentType: %+v | req.ReferrerID: %+v | req.ReferredID: %+v", req.PaymentType, req.ReferrerID, req.ReferralID)

	switch req.PaymentType {
	case referral_dto.PaymentPlatform:
		s.logger.Infof("req.ReferredID: %+v | req.ReferrerID: %+v | req.TicketCount: %+v", req.ReferralID, req.ReferrerID, req.TicketCount)

		bonusResult, err := s.calculateReferralBonuses(req, bonusRates, maxLevel)
		if err != nil {
			s.logger.Errorf("failed to calculate referral bonuses: %v", err)
			return "", errors.NewError(500, "failed to calculate referral bonuses")
		}

		jettonBalance, err := s.precheckoutBalance(s.config.PlatformSmartContract)
		if err != nil {
			s.logger.Errorf("failed to get jetton balance: %v", err)
			return "", errors.NewError(500, "failed to get jetton balance")
		}

		if jettonBalance < bonusResult.TotalBonusValue {
			s.logger.Errorf("insufficient balance in smart contract for bonus: %f", bonusResult.TotalBonusValue)
			return "", errors.NewError(400, "insufficient balance in smart contract")
		}

		s.logger.Infof("creating a cell for a transaction with the values of referral bonus accruals")
		cell, err := s.referral_helper.CellTransferJettonsFromPlatform(bonusResult.AccrualDictionary)
		if err != nil {
			s.logger.Errorf("failed to create cell: %v", err)
			return "", errors.NewError(500, "failed to create cell")
		}

		s.logger.Infof("transaction cell was created successfully: %+v", cell)

		s.logger.Infof("debiting balance for user %d: %f", req.ReferrerID, bonusResult.TotalBonusValue)
		s.logger.Infof("wallet seed: %v", len(s.config.WalletSeed))
		adminWallet, err := wallet.FromSeed(s.ton_client, s.config.WalletSeed, wallet.V4R2)
		if err != nil {
			s.logger.Errorf("failed to create wallet: %v", err)
			return "", errors.NewError(500, "failed to create wallet")
		}

		s.logger.Infof("wallet created successfully: %+v", adminWallet.Address())

		s.logger.Infof("sending a transaction to the smart contract")
		tx, _, err := adminWallet.SendWaitTransaction(context.Background(), &wallet.Message{
			Mode: wallet.PayGasSeparately,
			InternalMessage: &tlb.InternalMessage{
				Bounce:  true,
				DstAddr: address.MustParseAddr(s.config.PlatformSmartContract),
				Amount:  tlb.MustFromTON("0.1"),
				Body:    cell,
			},
		})

		if err != nil {
			s.logger.Errorf("transaction execution failed with an error: %v", err)
			return "", errors.NewError(500, "transaction execution failed")
		}

		s.logger.Info("transaction was completed successfully")
		s.logger.Infof("the hash of the transaction: %s", base64.StdEncoding.EncodeToString(tx.Hash))

		return "It is not intended for this type", nil
	case referral_dto.PaymentLeader:
		s.logger.Infof("req.ReferredID: %+v | req.ReferrerID: %+v | req.TicketCount: %+v | req.Amount: %+v", req.ReferralID, req.ReferrerID, req.TicketCount, req.AuthorID)
		if req.AuthorID == 0 {
			s.logger.Warnf("author ID is required for payment type %s", req.PaymentType)
			return "", errors.NewError(400, "author ID is required for payment type referred")
		}

		s.logger.Infof("fetching author data for user_id=%d", req.AuthorID)
		authorData, err := s.getAuthorData(req.AuthorID)
		if err != nil {
			s.logger.Errorf("failed to fetch author data: %v", err)
			return "", errors.NewError(500, "failed to fetch author data")
		}

		s.logger.Infof("author data fetched successfully: %+v", authorData)

		bonusResult, err := s.calculateReferralBonuses(req, bonusRates, maxLevel)
		if err != nil {
			s.logger.Errorf("failed to calculate referral bonuses: %v", err)
			return "", errors.NewError(500, "failed to calculate referral bonuses")
		}

		jettonBalance, err := s.precheckoutBalance(authorData.WalletAddress)
		if err != nil || jettonBalance < bonusResult.TotalBonusValue {
			s.logger.Errorf("failed to get jetton balance: or insufficient balance in author wallet for bonus: %v", err)
			s.logger.Infof("creating a payment order")

			newOrder := referral_adapters.CreatePaymentOrder(ctx, referral_dto.PaymentOrder{
				AuthorID:    req.AuthorID,
				ReferrerID:  req.ReferrerID,
				ReferralID:  req.ReferralID,
				TotalAmount: bonusResult.TotalBonusValue,
				TicketCount: req.TicketCount,
				CreatedAt:   time.Now().Unix(),
				Levels:      bonusResult.Levels,
			})

			err := s.referral_repository.CreatePaymentOrder(ctx, newOrder)
			if err != nil {
				s.logger.Errorf("failed to create payment order: %v", err)
				return "", errors.NewError(500, "failed to create payment order")
			}

			s.logger.Infof("payment order created successfully: %+v", newOrder)

			return "", errors.NewError(422, "failed to get jetton balance create payment order")
		}

		s.logger.Infof("creating a cell for a transaction with the values of referral bonus accruals")
		cell, err := s.referral_helper.CellTransferJettonsFromLeader(bonusResult.AccrualDictionary, bonusResult.TotalBonusValue)
		if err != nil {
			s.logger.Errorf("failed to create cell: %v", err)
			return "", errors.NewError(500, "failed to create cell")
		}

		s.logger.Infof("transaction cell was created successfully: %+v", cell)

		return base64.StdEncoding.EncodeToString(cell.ToBOC()), nil
	default:
		s.logger.Errorf("invalid payment type: %s", req.PaymentType)
		return "", errors.NewError(400, "invalid payment type")
	}
}
