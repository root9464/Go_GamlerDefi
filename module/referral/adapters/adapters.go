package referral_adapters

import (
	"context"

	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	referral_model "github.com/root9464/Go_GamlerDefi/module/referral/model"
)

func CreatePaymentOrder(ctx context.Context, req referral_dto.PaymentOrder) referral_model.PaymentOrder {
	levels := make([]referral_model.Level, len(req.Levels))
	for i, level := range req.Levels {
		levels[i] = referral_model.Level{
			LevelNumber: level.LevelNumber,
			Rate:        level.Rate,
			Amount:      level.Amount,
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
