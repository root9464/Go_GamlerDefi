package validation_repository

import (
	"context"

	validation_model "github.com/root9464/Go_GamlerDefi/src/module/validation/model"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IValidationRepository interface {
	CreateTransactionObserver(ctx context.Context, transaction validation_model.WorkerTransaction) (validation_model.WorkerTransaction, error)
	GetTransactionObserver(ctx context.Context, transactionID bson.ObjectID) (validation_model.WorkerTransaction, error)
	UpdateStatus(ctx context.Context, transactionID bson.ObjectID, status validation_model.WorkerStatus) (validation_model.WorkerTransaction, error)
	PrecheckoutTransaction(ctx context.Context, transactionID bson.ObjectID) (validation_model.WorkerTransaction, error)
	DeleteTransactionObserver(ctx context.Context, transactionID bson.ObjectID) error
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
