package main

import (
	"context"

	"github.com/root9464/Go_GamlerDefi/database"
	referral_repository "github.com/root9464/Go_GamlerDefi/module/referral/repository"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
)

const (
	database_url = "mongodb://root:example@localhost:27017"
)

func main() {
	logger := logger.GetLogger()
	client, err := database.ConnectDatabase(database_url, logger)
	if err != nil {
		logger.Error("❌ Failed to connect to MongoDB")
		return
	}

	referral_repository := referral_repository.NewReferralRepository(client, logger)

	// err = referral_repository.CreatePaymentOrder(context.Background(), referral_model.PaymentOrder{
	// 	AuthorID:    1,
	// 	ReferrerID:  2,
	// 	ReferralID:  3,
	// 	TotalAmount: 100,
	// 	TicketCount: 10,
	// 	CreatedAt:   time.Now(),
	// 	Levels: []referral_model.Level{
	// 		{LevelNumber: 1, Rate: 0.1, Amount: 10},
	// 	},
	// })

	// if err != nil {
	// 	logger.Error("❌ Failed to create payment order")
	// 	return
	// }

	// logger.Success("✅ Payment order created successfully")

	orders, err := referral_repository.GetPaymentOrdersByAuthorID(context.Background(), 1)
	if err != nil {
		logger.Error("❌ Failed to get payment orders")
		return
	}

	// Обработка результата
	if len(orders) == 0 {
		logger.Info("ℹ️ No payment orders found for AuthorID 1")
	} else {
		logger.Success("✅ Payment orders retrieved successfully")
		logger.Infof("orders: %+v", orders)
	}
}
