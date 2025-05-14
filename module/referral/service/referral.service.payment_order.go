package referral_service

import (
	"context"
	"encoding/base64"

	referral_adapters "github.com/root9464/Go_GamlerDefi/module/referral/adapters"
	referral_helper "github.com/root9464/Go_GamlerDefi/module/referral/helpers"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"github.com/xssnick/tonutils-go/address"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *ReferralService) PayPaymentOrder(ctx context.Context, paymentOrderID string) (string, error) {
	s.logger.Infof("start pay payment order: %s", paymentOrderID)
	orderID, err := bson.ObjectIDFromHex(paymentOrderID)
	if err != nil {
		s.logger.Errorf("failed to convert payment order ID to ObjectID: %v", err)
		return "", errors.NewError(500, "failed to convert payment order ID to ObjectID")
	}

	s.logger.Infof("fetching payment order in database by ID: %s", paymentOrderID)
	paymentOrder, err := s.referral_repository.GetPaymentOrderByID(ctx, orderID)
	if err != nil {
		s.logger.Errorf("failed to get payment order: %v", err)
		return "", errors.NewError(500, "failed to get payment order")
	}

	s.logger.Infof("payment order fetched successfully: %+v", paymentOrder)

	s.logger.Infof("converting payment order to DTO")
	paymentOrderDTO := referral_adapters.CreatePaymentOrderFromModel(ctx, paymentOrder)

	s.logger.Infof("converted payment order to DTO: %+v", paymentOrderDTO)

	s.logger.Infof("fetching author data for user_id=%d", paymentOrderDTO.AuthorID)
	authorData, err := s.getAuthorData(paymentOrderDTO.AuthorID)
	if err != nil {
		s.logger.Errorf("failed to get author data: %v", err)
		return "", errors.NewError(500, "failed to get author data")
	}

	s.logger.Infof("author data fetched successfully: %+v", authorData)

	s.logger.Infof("getting balance of author wallet")
	balance, err := s.precheckoutBalance(authorData.WalletAddress)
	if err != nil {
		s.logger.Errorf("failed to get balance of author wallet: %v", err)
		return "", errors.NewError(500, "failed to get balance of author wallet")
	}

	s.logger.Infof("balance of author wallet: %f", balance)

	if balance < paymentOrderDTO.TotalAmount {
		s.logger.Infof("insufficient funds on the balance sheet to pay the debt: %f", paymentOrderDTO.TotalAmount)
		return "", errors.NewError(402, "insufficient funds on the balance sheet to pay the debt")
	}

	s.logger.Infof("creating accrual dictionary for payment order")
	accrualDictionary := []referral_helper.JettonEntry{}
	for _, level := range paymentOrderDTO.Levels {
		accrualDictionary = append(accrualDictionary, referral_helper.JettonEntry{
			Address: address.MustParseAddr(level.Address),
			Amount:  uint64(level.Amount),
		})
	}
	s.logger.Infof("accrual dictionary created successfully: %+v", accrualDictionary)

	s.logger.Infof("creating a cell for a transaction with the values of referral bonus accruals")
	cell, err := s.referral_helper.CellTransferJettonsFromPlatform(accrualDictionary)
	if err != nil {
		s.logger.Errorf("failed to create cell: %v", err)
		return "", errors.NewError(500, "failed to create cell")
	}

	s.logger.Infof("transaction cell was created successfully: %+v", cell)

	return base64.StdEncoding.EncodeToString(cell.ToBOC()), nil
}
