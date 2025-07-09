package referral_repository

import (
	"context"

	referral_model "github.com/root9464/Go_GamlerDefi/module/referral/model"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IReferralRepository interface {
	CreatePaymentOrder(ctx context.Context, order referral_model.PaymentOrder) error
	GetPaymentOrderByID(ctx context.Context, orderID bson.ObjectID) (referral_model.PaymentOrder, error)
	GetPaymentOrdersByAuthorID(ctx context.Context, authorID int) ([]referral_model.PaymentOrder, error)
	GetAllPaymentOrders(ctx context.Context) ([]referral_model.PaymentOrder, error)
	DeleteAllPaymentOrders(ctx context.Context, authorID int) error
	DeletePaymentOrder(ctx context.Context, orderID bson.ObjectID) error
	GetDebtFromAuthorToReferrer(ctx context.Context, authorID int, referrerID int) ([]referral_model.PaymentOrder, error)
	UpdatePaymentOrder(ctx context.Context, order referral_model.PaymentOrder) error
	AddTrHashToPaymentOrder(ctx context.Context, orderID bson.ObjectID, trHash string) error
}

type ReferralRepository struct {
	logger *logger.Logger
	db     *mongo.Database
}

const (
	database_name             = "referral"
	payment_orders_collection = "payment_orders"
)

func NewReferralRepository(logger *logger.Logger, db *mongo.Database) IReferralRepository {
	return &ReferralRepository{logger: logger, db: db}
}
