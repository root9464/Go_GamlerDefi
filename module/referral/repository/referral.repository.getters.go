package referral_repository

import (
	"context"

	referral_model "github.com/root9464/Go_GamlerDefi/module/referral/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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

func (r *ReferralRepository) GetDebtFromAuthorToReferrer(ctx context.Context, authorID int, referrerID int) ([]referral_model.PaymentOrder, error) {
	r.logger.Info("getting payment orders by author ID and referrer ID")
	r.logger.Infof("author ID: %d, referrer ID: %d", authorID, referrerID)

	database := r.db.Database(database_name)
	collection := database.Collection(payment_orders_collection)

	filter := bson.D{
		{Key: "author_id", Value: authorID},
		{Key: "referrer_id", Value: referrerID},
	}

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

	r.logger.Infof("found %d payment orders for author ID %d and referrer ID %d", len(orders), authorID, referrerID)
	return orders, nil
}
