package referral_repository

import (
	"context"

	referral_model "github.com/root9464/Go_GamlerDefi/module/referral/model"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IReferralRepository interface {
	CreatePaymentOrder(ctx context.Context, order referral_model.PaymentOrder) error
	GetPaymentOrdersByAuthorID(ctx context.Context, authorID int) ([]referral_model.PaymentOrder, error)
	GetAllPaymentOrders(ctx context.Context) ([]referral_model.PaymentOrder, error)
}

type ReferralRepository struct {
	logger *logger.Logger
	db     *mongo.Client
}

const (
	database_name             = "referral"
	payment_orders_collection = "payment_orders"
)

func NewReferralRepository(db *mongo.Client, logger *logger.Logger) IReferralRepository {
	return &ReferralRepository{db: db, logger: logger}
}
