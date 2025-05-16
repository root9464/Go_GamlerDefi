package validation_repository

import (
	"context"

	validation_tr_model "github.com/root9464/Go_GamlerDefi/module/validation/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *ValidationRepository) GetTransactionObserver(transactionID bson.ObjectID) (validation_tr_model.WorkerTransaction, error) {
	r.logger.Infof("getting transaction observer: %v", transactionID)

	collection := r.db.Collection(collection_name)

	filter := bson.D{{Key: "_id", Value: transactionID}}

	var transaction validation_tr_model.WorkerTransaction
	err := collection.FindOne(context.Background(), filter).Decode(&transaction)
	if err != nil {
		r.logger.Errorf("failed to get transaction observer: %v", err)
		return validation_tr_model.WorkerTransaction{}, err
	}

	r.logger.Infof("transaction observer found: %v", transaction)
	r.logger.Info("transaction observer found successfully")
	return transaction, nil
}
