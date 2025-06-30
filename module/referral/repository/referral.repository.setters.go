package referral_repository

import (
	"context"
	"fmt"
	"time"

	referral_model "github.com/root9464/Go_GamlerDefi/module/referral/model"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (r *ReferralRepository) DeleteAllPaymentOrders(ctx context.Context) error {
	r.logger.Info("deleting all payment orders from database")

	collection := r.db.Collection(payment_orders_collection)

	result, err := collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		r.logger.Errorf("failed to delete all payment orders: %v", err)
		return err
	}

	r.logger.Infof("deleted %v payment orders", result.DeletedCount)
	return nil
}

type levelKey struct {
	LevelNumber int
	Address     string
}

func (r *ReferralRepository) UpdatePaymentOrder(ctx context.Context, order referral_model.PaymentOrder) error {
	r.logger.Info("updating payment order in database")
	r.logger.Infof("order: %+v", order)

	collection := r.db.Collection(payment_orders_collection)

	filter := bson.D{
		{Key: "leader_id", Value: order.LeaderID},
		{Key: "referrer_id", Value: order.ReferrerID},
		{Key: "referral_id", Value: order.ReferralID},
	}

	var existing referral_model.PaymentOrder
	err := collection.FindOne(ctx, filter).Decode(&existing)
	if err == mongo.ErrNoDocuments {
		r.logger.Errorf("payment order not found for leaderID %v", order.LeaderID)
		return err
	}
	if err != nil {
		r.logger.Errorf("failed to find payment order: %v", err)
		return err
	}

	mergedLevels := lo.Reduce(append(existing.Levels, order.Levels...), func(acc []referral_model.Level, level referral_model.Level, _ int) []referral_model.Level {
		key := levelKey{LevelNumber: level.LevelNumber, Address: level.Address}
		_, index, found := lo.FindIndexOf(acc, func(l referral_model.Level) bool {
			return l.LevelNumber == key.LevelNumber && l.Address == key.Address
		})

		if found {
			existingAmt, _ := decimal.NewFromString(acc[index].Amount.String())
			newAmt, _ := decimal.NewFromString(level.Amount.String())
			sum := existingAmt.Add(newAmt)
			sumDecimal, _ := bson.ParseDecimal128(sum.String())
			acc[index].Amount = sumDecimal
			return acc
		}
		return append(acc, level)
	}, existing.Levels[:0])

	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "total_amount", Value: order.TotalAmount},
			{Key: "ticket_count", Value: order.TicketCount},
		}},
		{Key: "$set", Value: bson.D{
			{Key: "levels", Value: mergedLevels},
			{Key: "updated_at", Value: time.Now().Unix()},
		}},
	}

	var updatedDoc referral_model.PaymentOrder
	err = collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedDoc)
	if err != nil {
		r.logger.Errorf("failed to update payment order: %v", err)
		return err
	}

	r.logger.Infof("payment order with leaderID %v updated successfully", order.LeaderID)
	r.logger.Infof("updated payment order: %+v", updatedDoc)
	return nil
}
