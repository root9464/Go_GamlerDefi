package referral_adapters

import (
	"fmt"

	referral_dto "github.com/root9464/Go_GamlerDefi/src/modules/referral/dto"
	referral_model "github.com/root9464/Go_GamlerDefi/src/modules/referral/model"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreatePaymentOrderFromDTO(req referral_dto.PaymentOrder) (referral_model.PaymentOrder, error) {
	levels := make([]referral_model.Level, len(req.Levels))
	for i, level := range req.Levels {
		rate, err := bson.ParseDecimal128(level.Rate.String())
		if err != nil {
			return referral_model.PaymentOrder{}, fmt.Errorf("failed to convert rate: %w", err)
		}

		amount, err := bson.ParseDecimal128(level.Amount.String())
		if err != nil {
			return referral_model.PaymentOrder{}, fmt.Errorf("failed to convert amount: %w", err)
		}

		levels[i] = referral_model.Level{
			LevelNumber: level.LevelNumber,
			Rate:        rate,
			Amount:      amount,
			Address:     level.Address,
		}
	}

	totalAmount, err := bson.ParseDecimal128(req.TotalAmount.String())
	if err != nil {
		return referral_model.PaymentOrder{}, fmt.Errorf("failed to convert total amount: %w", err)
	}

	paymentOrder := referral_model.PaymentOrder{
		LeaderID:    req.LeaderID,
		ReferrerID:  req.ReferrerID,
		ReferralID:  req.ReferralID,
		TotalAmount: totalAmount,
		TicketCount: req.TicketCount,
		CreatedAt:   req.CreatedAt,
		Levels:      levels,
		TrHash:      req.TrHash,
	}

	return paymentOrder, nil
}

func CreatePaymentOrderFromModel(dbData referral_model.PaymentOrder) (referral_dto.PaymentOrder, error) {
	levels := make([]referral_dto.LevelRequest, len(dbData.Levels))
	for i, level := range dbData.Levels {
		rate, err := decimal.NewFromString(level.Rate.String())
		if err != nil {
			return referral_dto.PaymentOrder{}, fmt.Errorf("failed to convert rate: %w", err)
		}

		amount, err := decimal.NewFromString(level.Amount.String())
		if err != nil {
			return referral_dto.PaymentOrder{}, fmt.Errorf("failed to convert amount: %w", err)
		}

		levels[i] = referral_dto.LevelRequest{
			LevelNumber: level.LevelNumber,
			Rate:        rate,
			Amount:      amount,
			Address:     level.Address,
		}
	}

	totalAmount, err := decimal.NewFromString(dbData.TotalAmount.String())
	if err != nil {
		return referral_dto.PaymentOrder{}, fmt.Errorf("failed to convert total amount: %w", err)
	}

	paymentOrderDTO := referral_dto.PaymentOrder{
		ID:          dbData.ID.Hex(),
		LeaderID:    dbData.LeaderID,
		ReferrerID:  dbData.ReferrerID,
		ReferralID:  dbData.ReferralID,
		TotalAmount: totalAmount,
		TicketCount: dbData.TicketCount,
		CreatedAt:   dbData.CreatedAt,
		Levels:      levels,
		TrHash:      dbData.TrHash,
	}

	return paymentOrderDTO, nil
}

func CreatePaymentOrderFromModelList(req []referral_model.PaymentOrder) ([]referral_dto.PaymentOrder, error) {
	paymentOrders := make([]referral_dto.PaymentOrder, len(req))
	for i, order := range req {
		paymentOrder, err := CreatePaymentOrderFromModel(order)
		if err != nil {
			return nil, fmt.Errorf("failed to create payment order from model: %w", err)
		}
		paymentOrders[i] = paymentOrder
	}

	return paymentOrders, nil
}
