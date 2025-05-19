package validation_repository

import (
	validation_model "github.com/root9464/Go_GamlerDefi/module/validation/model"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IValidationRepository interface {
	CreateTransactionObserver(transaction validation_model.WorkerTransaction) (validation_model.WorkerTransaction, error)
	GetTransactionObserver(transactionID bson.ObjectID) (validation_model.WorkerTransaction, error)
	UpdateStatus(transactionID bson.ObjectID, status validation_model.WorkerStatus) (validation_model.WorkerTransaction, error)
	PrecheckoutTransaction(transactionID bson.ObjectID) (validation_model.WorkerTransaction, error)
	DeleteTransactionObserver(transactionID bson.ObjectID) error
}

type ValidationRepository struct {
	logger *logger.Logger
	db     *mongo.Database
}

const (
	collection_name = "validation_transaction"
)

func NewValidationRepository(logger *logger.Logger, db *mongo.Database) IValidationRepository {
	return &ValidationRepository{logger: logger, db: db}
}
