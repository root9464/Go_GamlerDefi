package validation_repository

import (
	"context"
	"time"

	validation_model "github.com/root9464/Go_GamlerDefi/src/module/validation/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *ValidationRepository) CreateTransactionObserver(ctx context.Context, transaction validation_model.WorkerTransaction) (validation_model.WorkerTransaction, error) {
	r.logger.Infof("creating transaction observer: %v", transaction)

	if transaction.CreatedAt == 0 {
		r.logger.Infof("transaction.CreatedAt is zero, setting to current time")
		transaction.CreatedAt = time.Now().Unix()
	}

	if transaction.UpdatedAt == 0 {
		r.logger.Infof("transaction.UpdatedAt is zero, setting to current time")
		transaction.UpdatedAt = time.Now().Unix()
	}

	if transaction.ID.IsZero() {
		r.logger.Infof("transaction.ID is zero, setting to new ObjectID")
		transaction.ID = bson.NewObjectID()
	}

	collection := r.db.Collection(collection_name)

	result, err := collection.InsertOne(ctx, transaction)
	if err != nil {
		r.logger.Errorf("failed to insert transaction observer: %v", err)
		return validation_model.WorkerTransaction{}, err
	}

	r.logger.Infof("transaction observer created with ID: %v", result.InsertedID)
	r.logger.Infof("transaction observer created successfully: %v", transaction)
	return transaction, nil
}

func (r *ValidationRepository) UpdateStatus(ctx context.Context, transactionID bson.ObjectID, status validation_model.WorkerStatus) (validation_model.WorkerTransaction, error) {
	r.logger.Infof("updating status for transaction: %v", transactionID)

	collection := r.db.Collection(collection_name)

	filter := bson.D{{Key: "_id", Value: transactionID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: status},
			{Key: "updated_at", Value: time.Now().Unix()},
		}},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedDoc validation_model.WorkerTransaction
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDoc)

	if err != nil {
		r.logger.Errorf("failed to update status: %v", err)
		return validation_model.WorkerTransaction{}, err
	}

	r.logger.Infof("status updated for transaction: %v", transactionID)
	r.logger.Infof("new status order transaction: %v", updatedDoc.Status)
	return updatedDoc, nil
}

func (r *ValidationRepository) PrecheckoutTransaction(ctx context.Context, transactionID bson.ObjectID) (validation_model.WorkerTransaction, error) {
	r.logger.Infof("get transaction in db: %+v", transactionID)

	transactionObserver, err := r.GetTransactionObserver(ctx, transactionID)
	if err != nil {
		r.logger.Errorf("failed to get transaction from database: %v", err)
		return validation_model.WorkerTransaction{}, err
	}

	r.logger.Infof("current transaction status: %v (%T)", transactionObserver.Status, transactionObserver.Status)
	if transactionObserver.Status != validation_model.WorkerStatusRunning {
		r.logger.Warnf("transaction status is not running: %v", transactionObserver.Status)
	}

	r.logger.Info("update transaction status to running")
	transaction, err := r.UpdateStatus(ctx, transactionID, validation_model.WorkerStatusRunning)
	if err != nil {
		r.logger.Errorf("failed to update status: %v", err)
		return validation_model.WorkerTransaction{}, err
	}

	r.logger.Infof("transaction status updated to running: %v", transaction.Status)

	return transaction, nil
}

func (r *ValidationRepository) DeleteTransactionObserver(ctx context.Context, transactionID bson.ObjectID) error {
	r.logger.Infof("deleting transaction observer: %v", transactionID)

	collection := r.db.Collection(collection_name)

	_, err := collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: transactionID}})
	if err != nil {
		r.logger.Errorf("failed to delete transaction observer: %v", err)
		return err
	}

	r.logger.Infof("transaction observer deleted: %v", transactionID)
	return nil
}
