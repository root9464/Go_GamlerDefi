package referral_service

import (
	"context"
	"fmt"
	"iter"
	"math"
	"strconv"

	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	referral_helper "github.com/root9464/Go_GamlerDefi/module/referral/helpers"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"github.com/root9464/Go_GamlerDefi/packages/utils"
	"github.com/samber/lo"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

type ReferralLevel struct {
	Level         int
	Rate          float64
	ReferrerID    int
	WalletAddress string
	Err           error
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

func (s *ReferralService) ReferralChainIterator(req referral_dto.ReferralProcessRequest, bonusRates map[int]float64, maxLevel int) iter.Seq[ReferralLevel] {
	return func(yield func(ReferralLevel) bool) {
		currentReferrerID := req.ReferrerID
		referredID := req.ReferredID

		s.logger.Infof("fetching first level referrer for user %d", currentReferrerID)
		referrerL1, err := s.GetReferrerChain(currentReferrerID)
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

			referrerData, err := s.GetReferrerChain(currentReferrerID)
			if err != nil {
				s.logger.Errorf("failed to fetch referrer data for user %d at level %d: %v", currentReferrerID, level, err)
				yield(ReferralLevel{Level: level, Rate: rate, ReferrerID: currentReferrerID, WalletAddress: "", Err: err})
				return
			}

			if !yield(ReferralLevel{Level: level, Rate: rate, ReferrerID: currentReferrerID, WalletAddress: referrerData.WalletAddress, Err: nil}) {
				s.logger.Infof("stopping referral chain at level %d", level+1)
				return
			}

			parentData, err := s.GetReferrerChain(currentReferrerID)
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

func (s *ReferralService) CalculateReferralBonuses(ctx context.Context, req referral_dto.ReferralProcessRequest) (string, error) {
	s.logger.Infof("starting referral bonus calculation for: %+v", req)

	bonusRates := map[int]float64{
		0: 0.20, // Уровень 1: 20%
		1: 0.02, // Уровень 2: 2%
	}

	s.logger.Infof("bonusRates: %+v", bonusRates)
	s.logger.Infof("req.PaymentType: %+v | req.ReferrerID: %+v | req.ReferredID: %+v", req.PaymentType, req.ReferrerID, req.ReferredID)

	switch req.PaymentType {
	case referral_dto.PaymentAuthor:
		s.logger.Infof("req.ReferredID: %+v | req.ReferrerID: %+v | req.TicketCount: %+v", req.ReferredID, req.ReferrerID, req.TicketCount)
		s.logger.Infof("accruing bonus for levels maxLevel: %d", maxLevel)

		totalBonusValue := 0.0
		accrualDictionary := []referral_helper.JettonEntry{}

		for referralLevel := range s.ReferralChainIterator(req, bonusRates, maxLevel) {
			if referralLevel.Err != nil {
				s.logger.Errorf("error in referral chain at level %d: %v", referralLevel.Level, referralLevel.Err)
				return "", referralLevel.Err
			}

			rate := referralLevel.Rate
			bonusAmount := (math.Round(float64(req.TicketCount)*rate*100) / 100)
			totalBonusValue += bonusAmount
			if referralLevel.WalletAddress != "" {
				accrualDictionary = append(accrualDictionary, referral_helper.JettonEntry{
					Address: address.MustParseAddr(referralLevel.WalletAddress),
					Amount:  uint64(bonusAmount * math.Pow10(9)),
				})
			}

			s.logger.Infof("level %d: referrer %d (%s) gets %.2f bonus", referralLevel.Level, referralLevel.ReferrerID, referralLevel.WalletAddress, bonusAmount)
		}

		s.logger.Infof("total bonus value: %f", totalBonusValue)
		s.logger.Infof("dictionary with the values of referral bonus accruals: %+v", accrualDictionary)

		s.logger.Infof("checking the balance of a smart contract for awarding bonuses")
		contractBalance, err := s.ton_api.GetAccountJettonsBalances(context.Background(), tonapi.GetAccountJettonsBalancesParams{
			AccountID: s.config.PlatformSmartContract,
		})

		if err != nil {
			s.logger.Errorf("failed to fetch account jettons balances: %v", err)
			return "", errors.NewError(500, "failed to fetch account jettons balances")
		}

		s.logger.Infof("find smart contract address %s in balances", s.config.PlatformSmartContract)
		foundJetton, found := lo.Find(contractBalance.Balances, func(b tonapi.JettonBalance) bool {
			rawAddr, parseErr := address.ParseRawAddr(b.WalletAddress.Address)
			if parseErr != nil {
				s.logger.Errorf("failed to parse wallet address: %v", parseErr)
				return false
			}
			userFriendlyAddr := rawAddr.Bounce(true).Testnet(true).String()
			s.logger.Infof("user friendly address: %s", userFriendlyAddr)
			return userFriendlyAddr == s.config.SmartContractJettonWallet
		})
		if !found {
			s.logger.Errorf("Smart contract address %s not found in balances", s.config.PlatformSmartContract)
			return "", errors.NewError(404, "smart contract address not found")
		}
		s.logger.Infof("Smart contract address %s found in balances %s", s.config.PlatformSmartContract, foundJetton.Balance)

		jettonBalance, err := strconv.ParseFloat(foundJetton.Balance, 64)
		if err != nil {
			s.logger.Errorf("failed to convert jetton balance to float64: %v", err)
			return "", errors.NewError(500, "failed to convert jetton balance to float64")
		}

		if jettonBalance/math.Pow10(foundJetton.Jetton.Decimals) < totalBonusValue {
			s.logger.Errorf("insufficient balance in smart contract for bonus: %f", totalBonusValue)
			return "", errors.NewError(400, "insufficient balance in smart contract")
		}

		s.logger.Infof("creating a cell for a transaction with the values of referral bonus accruals")
		cell, err := s.referral_helper.CellTransferJettonsFromPlatform(accrualDictionary)
		if err != nil {
			s.logger.Errorf("failed to create cell: %v", err)
			return "", errors.NewError(500, "failed to create cell")
		}

		s.logger.Infof("transaction cell was created successfully: %+v", cell)

		s.logger.Infof("debiting balance for user %d: %f", req.ReferrerID, totalBonusValue)
		s.logger.Infof("wallet seed: %v", len(s.config.WalletSeed))
		adminWallet, err := wallet.FromSeed(s.ton_client, s.config.WalletSeed, wallet.V4R2)
		if err != nil {
			s.logger.Errorf("failed to create wallet: %v", err)
			return "", errors.NewError(500, "failed to create wallet")
		}

		s.logger.Infof("wallet created successfully: %+v", adminWallet.Address())
		// логика вызова смарта и начисления
		return cell, nil
	case referral_dto.PaymentReferred:
		return "", nil

	default:
		s.logger.Errorf("invalid payment type: %s", req.PaymentType)
		return "", errors.NewError(400, "invalid payment type")
	}
}
