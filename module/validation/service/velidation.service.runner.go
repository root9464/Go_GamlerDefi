package validation_service

import (
	validation_adapters "github.com/root9464/Go_GamlerDefi/module/validation/adapters"
	validation_dto "github.com/root9464/Go_GamlerDefi/module/validation/dto"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *ValidationService) RunnerTransaction(transaction *validation_dto.WorkerTransactionDTO) (*validation_dto.WorkerTransactionDTO, bool, error) {
	s.logger.Info("start running validation transaction")

	s.logger.Info("validate dto")
	if err := s.validator.Struct(transaction); err != nil {
		s.logger.Errorf("failed to validate transaction dto: %v", err)
		return nil, false, err
	}
	s.logger.Info("validate dto success")

	if transaction.Status != validation_dto.WorkerStatusPending {
		s.logger.Errorf("transaction status is not pending")
		return nil, false, errors.NewError(400, "transaction status is not pending")
	}

	s.logger.Info("convert transaction id to bson.ObjectID")
	transactionID, err := bson.ObjectIDFromHex(transaction.ID)
	if err != nil {
		s.logger.Errorf("failed to convert transaction id to bson.ObjectID: %v", err)
		return nil, false, err
	}

	s.logger.Infof("get transaction in db: %+v", transactionID)

	transactionObserver, err := s.validation_repository.GetTransactionObserver(transactionID)
	switch err {
	case nil:
		s.logger.Info("transaction already exists in the database")
		s.logger.Info("convert transaction model to dto")
		transactionDTO := validation_adapters.TransactionModelToDTO(transactionObserver)
		s.logger.Info("transaction dto created successfully")
		return &transactionDTO, true, nil
	default:
		s.logger.Errorf("failed to get transaction from database: %v", err)

		s.logger.Infof("transaction observer: %+v", transactionObserver)
		if transactionObserver.ID != bson.NilObjectID {
			s.logger.Errorf("transaction already exists in database, RunValidation is not intended for existing transactions")
			return nil, false, errors.NewError(422, "transaction already exists in database, RunValidation is not intended for existing transactions")
		}

		transactionModel, err := validation_adapters.TransactionDTOToModel(*transaction)
		if err != nil {
			s.logger.Errorf("failed to convert transaction dto to model: %v", err)
			return nil, false, err
		}

		s.logger.Info("create transaction observer")
		transactionModel, err = s.validation_repository.CreateTransactionObserver(transactionModel)
		if err != nil {
			s.logger.Errorf("failed to create transaction observer: %v", err)
			return nil, false, err
		}

		s.logger.Infof("transaction observer created successfully: %+v", transactionModel.ID)
		s.logger.Info("convert transaction model to dto")
		transactionDTO := validation_adapters.TransactionModelToDTO(transactionModel)
		s.logger.Info("transaction dto created successfully")
		return &transactionDTO, true, nil
	}
}
