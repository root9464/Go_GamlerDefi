package validation_tr_repository

import (
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IValidationTrRepository interface {
}

type ValidationTrRepository struct {
	logger *logger.Logger
	db     *mongo.Client
}

const (
	database_name   = "validation"
	collection_name = "validation_tr"
)

func NewValidationTrRepository(logger *logger.Logger, db *mongo.Client) IValidationTrRepository {
	return &ValidationTrRepository{logger: logger, db: db}
}
