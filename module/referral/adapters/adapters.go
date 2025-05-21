package referral_adapters

import (
	"context"

	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	referral_model "github.com/root9464/Go_GamlerDefi/module/referral/model"
)

func CreatePaymentOrderFromDTO(ctx context.Context, req referral_dto.PaymentOrder) referral_model.PaymentOrder {
	levels := make([]referral_model.Level, len(req.Levels))
	for i, level := range req.Levels {
		levels[i] = referral_model.Level{
			LevelNumber: level.LevelNumber,
			Rate:        level.Rate,
			Amount:      level.Amount,
			Address:     level.Address,
		}
	}

	paymentOrder := referral_model.PaymentOrder{
		AuthorID:    req.AuthorID,
		ReferrerID:  req.ReferrerID,
		ReferralID:  req.ReferralID,
		TotalAmount: req.TotalAmount,
		TicketCount: req.TicketCount,
		CreatedAt:   req.CreatedAt,
		Levels:      levels,
	}

	return paymentOrder
}

func CreatePaymentOrderFromModel(ctx context.Context, dbData referral_model.PaymentOrder) referral_dto.PaymentOrder {
	levels := make([]referral_dto.LevelRequest, len(dbData.Levels))
	for i, level := range dbData.Levels {
		levels[i] = referral_dto.LevelRequest{
			LevelNumber: level.LevelNumber,
			Rate:        level.Rate,
			Amount:      level.Amount,
			Address:     level.Address,
		}
	}

	paymentOrderDTO := referral_dto.PaymentOrder{
		ID:          dbData.ID.Hex(),
		AuthorID:    dbData.AuthorID,
		ReferrerID:  dbData.ReferrerID,
		ReferralID:  dbData.ReferralID,
		TotalAmount: dbData.TotalAmount,
		TicketCount: dbData.TicketCount,
		CreatedAt:   dbData.CreatedAt,
		Levels:      levels,
	}

	return paymentOrderDTO
}

func CreatePaymentOrderFromModelList(ctx context.Context, req []referral_model.PaymentOrder) []referral_dto.PaymentOrder {
	paymentOrders := make([]referral_dto.PaymentOrder, len(req))
	for i, order := range req {
		paymentOrders[i] = CreatePaymentOrderFromModel(ctx, order)
	}

	return paymentOrders
}
