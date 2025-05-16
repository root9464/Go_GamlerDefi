package validation_repository

import (
	"context"
	"fmt"
	"time"

	validation_model "github.com/root9464/Go_GamlerDefi/module/validation/model"
	validation_tr_model "github.com/root9464/Go_GamlerDefi/module/validation/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *ValidationRepository) CreateTransactionObserver(transaction validation_tr_model.WorkerTransaction) (validation_tr_model.WorkerTransaction, error) {
	r.logger.Infof("creating transaction observer: %v", transaction)

	if transaction.CreatedAt == 0 {
		r.logger.Infof("transaction.CreatedAt is zero, setting to current time")
		transaction.CreatedAt = time.Now().Unix()
	}

	if transaction.ID.IsZero() {
		r.logger.Infof("transaction.ID is zero, setting to new ObjectID")
		transaction.ID = bson.NewObjectID()
	}

	collection := r.db.Collection(collection_name)

	result, err := collection.InsertOne(context.Background(), transaction)
	if err != nil {
		r.logger.Errorf("failed to insert transaction observer: %v", err)
		return validation_tr_model.WorkerTransaction{}, err
	}

	r.logger.Infof("transaction observer created with ID: %v", result.InsertedID)
	r.logger.Infof("transaction observer created successfully: %v", transaction)
	return transaction, nil
}

func (r *ValidationRepository) UpdateStatus(transactionID bson.ObjectID, status validation_model.WorkerStatus) error {
	r.logger.Infof("updating status for transaction: %v", transactionID)

	collection := r.db.Collection(collection_name)

	filter := bson.D{{Key: "_id", Value: transactionID}}
	update := bson.D{{Key: "status", Value: status}, {Key: "updated_at", Value: time.Now().Unix()}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		r.logger.Errorf("failed to update status: %v", err)
		return err
	}

	r.logger.Infof("status updated for transaction: %v", transactionID)
	r.logger.Infof("new status order transaction: %v", status)
	return nil
}

func (r *ValidationRepository) PrecheckoutTransaction(transactionID bson.ObjectID) error {
	r.logger.Infof("get transaction in db: %+v", transactionID)

	transactionObserver, err := r.GetTransactionObserver(transactionID)
	if err != nil {
		r.logger.Errorf("failed to get transaction from database: %v", err)
		return err
	}

	if transactionObserver.Status != validation_model.WorkerStatusPending {
		r.logger.Errorf("transaction status is not pending")
		return fmt.Errorf("transaction status is not pending")
	}

	r.logger.Info("update transaction status to running")
	err = r.UpdateStatus(transactionID, validation_model.WorkerStatusRunning)
	if err != nil {
		r.logger.Errorf("failed to update status: %v", err)
		return err
	}

	return nil
}

func (r *ValidationRepository) DeleteTransactionObserver(transactionID bson.ObjectID) error {
	r.logger.Infof("deleting transaction observer: %v", transactionID)

	collection := r.db.Collection(collection_name)

	_, err := collection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: transactionID}})
	if err != nil {
		r.logger.Errorf("failed to delete transaction observer: %v", err)
		return err
	}

	r.logger.Infof("transaction observer deleted: %v", transactionID)
	return nil
}
