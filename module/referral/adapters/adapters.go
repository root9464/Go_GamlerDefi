package referral_adapters

import (
	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	referral_model "github.com/root9464/Go_GamlerDefi/module/referral/model"
)

func CreatePaymentOrderFromDTO(req referral_dto.PaymentOrder) referral_model.PaymentOrder {
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
		LeaderID:    req.LeaderID,
		ReferrerID:  req.ReferrerID,
		ReferralID:  req.ReferralID,
		TotalAmount: req.TotalAmount,
		TicketCount: req.TicketCount,
		CreatedAt:   req.CreatedAt,
		Levels:      levels,
	}

	return paymentOrder
}

func CreatePaymentOrderFromModel(dbData referral_model.PaymentOrder) referral_dto.PaymentOrder {
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
		LeaderID:    dbData.LeaderID,
		ReferrerID:  dbData.ReferrerID,
		ReferralID:  dbData.ReferralID,
		TotalAmount: dbData.TotalAmount,
		TicketCount: dbData.TicketCount,
		CreatedAt:   dbData.CreatedAt,
		Levels:      levels,
	}

	return paymentOrderDTO
}

func CreatePaymentOrderFromModelList(req []referral_model.PaymentOrder) []referral_dto.PaymentOrder {
	paymentOrders := make([]referral_dto.PaymentOrder, len(req))
	for i, order := range req {
		paymentOrders[i] = CreatePaymentOrderFromModel(order)
	}

	return paymentOrders
}
