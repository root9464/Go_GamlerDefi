package referral_repository

import (
	"context"
	"fmt"
	"time"

	referral_model "github.com/root9464/Go_GamlerDefi/module/referral/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *ReferralRepository) CreatePaymentOrder(ctx context.Context, order referral_model.PaymentOrder) error {
	r.logger.Info("create payment order in database")

	if order.CreatedAt == 0 {
		r.logger.Infof("order.CreatedAt is zero, setting to current time")
		order.CreatedAt = time.Now().Unix()
	}

	if order.ID.IsZero() {
		r.logger.Infof("order.ID is zero, setting to new ObjectID")
		order.ID = bson.NewObjectID()
	}

	collection := r.db.Collection(payment_orders_collection)

	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		r.logger.Errorf("failed to insert payment order: %v", err)
		return err
	}

	r.logger.Infof("payment order inserted with ID: %v", result)
	r.logger.Info("payment order created successfully")
	r.logger.Infof("created payment order: %+v", order)
	return nil
}

func (r *ReferralRepository) DeletePaymentOrder(ctx context.Context, orderID bson.ObjectID) error {
	r.logger.Info("deleting payment order from database")
	r.logger.Infof("order ID: %v", orderID)

	collection := r.db.Collection(payment_orders_collection)

	filter := bson.D{{Key: "_id", Value: orderID}}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Errorf("failed to delete payment order: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		r.logger.Warnf("no payment order found with ID: %v", orderID)
		return fmt.Errorf("payment order not found")
	}

	r.logger.Infof("payment order with ID %v deleted successfully", orderID)
	return nil
}
