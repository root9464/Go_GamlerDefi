package referral_repository

import (
	"context"
	"time"

	referral_model "github.com/root9464/Go_GamlerDefi/module/referral/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *ReferralRepository) CreatePaymentOrder(ctx context.Context, order referral_model.PaymentOrder) error {
	r.logger.Info("create payment order in database")
	r.logger.Infof("order data: %+v", order)

	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}

	database := r.db.Database(database_name)
	collection := database.Collection(payment_orders_collection)

	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		r.logger.Errorf("failed to insert payment order: %v", err)
		return err
	}

	r.logger.Infof("payment order inserted with ID: %v", result.InsertedID)
	r.logger.Info("payment order created successfully")
	return nil
}

func (r *ReferralRepository) GetPaymentOrdersByAuthorID(ctx context.Context, authorID int) ([]referral_model.PaymentOrder, error) {
	r.logger.Info("getting payment orders by author ID")
	r.logger.Infof("author ID: %d", authorID)

	database := r.db.Database(database_name)
	collection := database.Collection(payment_orders_collection)

	filter := bson.D{{Key: "author_id", Value: authorID}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		r.logger.Errorf("failed to find payment orders: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []referral_model.PaymentOrder
	if err := cursor.All(ctx, &orders); err != nil {
		r.logger.Errorf("failed to decode payment orders: %v", err)
		return nil, err
	}

	r.logger.Infof("found %d payment orders for author ID %d", len(orders), authorID)
	return orders, nil
}

func (r *ReferralRepository) GetAllPaymentOrders(ctx context.Context) ([]referral_model.PaymentOrder, error) {
	r.logger.Info("Getting all payment orders")
	database := r.db.Database(database_name)
	collection := database.Collection(payment_orders_collection)
	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		r.logger.Errorf("Failed to find all payment orders: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var orders []referral_model.PaymentOrder
	if err := cursor.All(ctx, &orders); err != nil {
		r.logger.Errorf("Failed to decode all payment orders: %v", err)
		return nil, err
	}
	r.logger.Infof("Found %d payment orders", len(orders))
	return orders, nil
}
