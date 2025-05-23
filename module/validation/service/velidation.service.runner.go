package validation_service

import (
	"context"

	validation_adapters "github.com/root9464/Go_GamlerDefi/module/validation/adapters"
	validation_dto "github.com/root9464/Go_GamlerDefi/module/validation/dto"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *ValidationService) RunnerTransaction(ctx context.Context, transaction *validation_dto.WorkerTransactionDTO) (*validation_dto.WorkerTransactionDTO, bool, error) {
	s.logger.Info("start running validation transaction")
	if transaction.Status != validation_dto.WorkerStatusPending {
		s.logger.Errorf("transaction status is not pending")
		return transaction, false, errors.NewError(400, "transaction status is not pending")
	}

	if transaction.ID == "" {
		s.logger.Warnf("transaction id is empty, creating new transaction id")
		transaction.ID = bson.NewObjectID().Hex()
	}

	s.logger.Info("convert transaction id to bson.ObjectID")
	transactionID, err := bson.ObjectIDFromHex(transaction.ID)
	if err != nil {
		s.logger.Errorf("failed to convert transaction id to bson.ObjectID: %v", err)
		return transaction, false, err
	}

	s.logger.Infof("get transaction in db: %+v", transactionID)

	transactionObserver, err := s.validation_repository.GetTransactionObserver(ctx, transactionID)
	switch err {
	case nil:
		s.logger.Info("transaction already exists in the database")
		s.logger.Info("convert transaction model to dto")
		transactionDTO := validation_adapters.TransactionModelToDTOPoint(transactionObserver)
		s.logger.Info("transaction dto created successfully")
		return transactionDTO, true, nil
	default:
		s.logger.Errorf("failed to get transaction from database: %v", err)
		if transactionObserver.ID != bson.NilObjectID {
			s.logger.Errorf("transaction already exists in database, RunValidation is not intended for existing transactions")
			return transaction, false, errors.NewError(422, "transaction already exists in database, RunValidation is not intended for existing transactions")
		}

		transactionModel, err := validation_adapters.TransactionDTOToModel(*transaction)
		if err != nil {
			s.logger.Errorf("failed to convert transaction dto to model: %v", err)
			return transaction, false, err
		}

		s.logger.Info("create transaction observer")
		transactionModel, err = s.validation_repository.CreateTransactionObserver(ctx, transactionModel)
		if err != nil {
			s.logger.Errorf("failed to create transaction observer: %v", err)
			return transaction, false, err
		}

		s.logger.Infof("transaction observer created successfully: %+v", transactionModel.ID)
		s.logger.Info("convert transaction model to dto")
		transactionDTO := validation_adapters.TransactionModelToDTOPoint(transactionModel)
		s.logger.Info("transaction dto created successfully")
		return transactionDTO, true, nil
	}
}
