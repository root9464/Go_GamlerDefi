package referral_service

import (
	"context"
	"encoding/base64"

	referral_adapters "github.com/root9464/Go_GamlerDefi/src/module/referral/adapters"
	referral_helper "github.com/root9464/Go_GamlerDefi/src/module/referral/helpers"
	errors "github.com/root9464/Go_GamlerDefi/src/packages/lib/error"
	"github.com/shopspring/decimal"
	"github.com/xssnick/tonutils-go/address"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *ReferralService) PayPaymentOrder(ctx context.Context, paymentOrderID string, walletAddress string) (string, error) {
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
	paymentOrderDTO, err := referral_adapters.CreatePaymentOrderFromModel(paymentOrder)
	if err != nil {
		s.logger.Errorf("failed to convert payment order to DTO: %v", err)
		return "", errors.NewError(500, "failed to convert payment order to DTO")
	}

	s.logger.Infof("converted payment order to DTO: %+v", paymentOrderDTO)

	jettonBalance, err := s.precheckoutBalance(s.config.PlatformSmartContract)
	if err != nil {
		s.logger.Errorf("failed to get jetton balance: %v", err)
		return "", errors.NewError(500, "failed to get jetton balance")
	}
	s.logger.Infof("jetton balance: %s", jettonBalance.String())

	if jettonBalance.LessThan(paymentOrderDTO.TotalAmount) {
		s.logger.Errorf("insufficient balance in smart contract for bonus: %s", paymentOrderDTO.TotalAmount.String())
		return "", errors.NewError(400, "insufficient balance in smart contract")
	}

	s.logger.Infof("fetching author data for user_id=%d", paymentOrderDTO.LeaderID)
	authorData, err := s.getAuthorData(paymentOrderDTO.LeaderID)
	if err != nil {
		s.logger.Errorf("failed to get author data: %v", err)
		return "", errors.NewError(500, "failed to get author data")
	}

	s.logger.Infof("author data fetched successfully: %+v", authorData)

	s.logger.Infof("getting balance of author wallet")
	balance, err := s.precheckoutBalance(walletAddress)
	if err != nil {
		s.logger.Errorf("failed to get balance of author wallet: %v", err)
		return "", errors.NewError(500, "failed to get balance of author wallet")
	}

	s.logger.Infof("balance of author wallet: %s", balance.String())

	if balance.LessThan(paymentOrderDTO.TotalAmount) {
		s.logger.Infof("insufficient funds on the balance sheet to pay the debt: %s", paymentOrderDTO.TotalAmount.String())
		return "", errors.NewError(402, "insufficient funds on the balance sheet to pay the debt")
	}

	s.logger.Infof("creating accrual dictionary for payment order")
	accrualDictionary := []referral_helper.JettonEntry{}
	for _, level := range paymentOrderDTO.Levels {
		accrualDictionary = append(accrualDictionary, referral_helper.JettonEntry{
			Address: address.MustParseAddr(level.Address),
			Amount:  level.Amount,
		})
	}
	s.logger.Infof("accrual dictionary created successfully: %+v", accrualDictionary)

	s.logger.Infof("creating a cell for a transaction with the values of referral bonus accruals")
	cell, err := s.referral_helper.CellTransferJettonsFromLeader(accrualDictionary, paymentOrderDTO.TotalAmount)
	if err != nil {
		s.logger.Errorf("failed to create cell: %v", err)
		return "", errors.NewError(500, "failed to create cell")
	}

	s.logger.Infof("transaction cell was created successfully: %+v", cell)

	return base64.StdEncoding.EncodeToString(cell.ToBOC()), nil
}

func (s *ReferralService) PayAllPaymentOrders(ctx context.Context, authorID int, walletAddress string) (string, error) {
	s.logger.Infof("start pay all payment orders for user_id=%d", authorID)

	s.logger.Infof("fetching payment orders in database by author_id: %d", authorID)
	paymentOrders, err := s.referral_repository.GetPaymentOrdersByAuthorID(ctx, authorID)
	if err != nil {
		s.logger.Errorf("failed to get payment orders: %v", err)
		return "", errors.NewError(500, "failed to get payment orders")
	}

	s.logger.Infof("payment orders fetched successfully: %+v", paymentOrders)

	s.logger.Infof("converting payment order to DTO")
	paymentOrderDTO, err := referral_adapters.CreatePaymentOrderFromModelList(paymentOrders)
	if err != nil {
		s.logger.Errorf("failed to convert payment order to DTO: %v", err)
		return "", errors.NewError(500, "failed to convert payment order to DTO")
	}

	s.logger.Infof("converted payment order to DTO: %+v", paymentOrderDTO)

	s.logger.Infof("fetching author data for user_id=%d", authorID)
	authorData, err := s.getAuthorData(authorID)
	if err != nil {
		s.logger.Errorf("failed to get author data: %v", err)
		return "", errors.NewError(500, "failed to get author data")
	}

	s.logger.Infof("author data fetched successfully: %+v", authorData)

	s.logger.Infof("getting balance of author wallet")
	balance, err := s.precheckoutBalance(walletAddress)
	if err != nil {
		s.logger.Errorf("failed to get balance of author wallet: %v", err)
		return "", errors.NewError(500, "failed to get balance of author wallet")
	}

	s.logger.Infof("balance of author wallet: %s", balance.String())

	totalAmount := decimal.NewFromFloat(0)
	for _, paymentOrder := range paymentOrderDTO {
		totalAmount = totalAmount.Add(paymentOrder.TotalAmount)
	}

	if balance.LessThan(totalAmount) {
		s.logger.Infof("insufficient funds on the balance sheet to pay the debt: %s", totalAmount.String())
		return "", errors.NewError(402, "insufficient funds on the balance sheet to pay the debt")
	}

	s.logger.Infof("creating accrual dictionary for payment order")
	accrualDictionary := []referral_helper.JettonEntry{}
	for _, paymentOrder := range paymentOrderDTO {
		for _, level := range paymentOrder.Levels {
			accrualDictionary = append(accrualDictionary, referral_helper.JettonEntry{
				Address: address.MustParseAddr(level.Address),
				Amount:  level.Amount,
			})
		}
	}

	s.logger.Infof("accrual dictionary created successfully: %+v", accrualDictionary)

	s.logger.Infof("creating a cell for a transaction with the values of referral bonus accruals")
	cell, err := s.referral_helper.CellTransferJettonsFromLeader(accrualDictionary, totalAmount)
	if err != nil {
		s.logger.Errorf("failed to create cell: %v", err)
		return "", errors.NewError(500, "failed to create cell")
	}

	s.logger.Infof("transaction cell was created successfully: %+v", cell)

	return base64.StdEncoding.EncodeToString(cell.ToBOC()), nil
}
