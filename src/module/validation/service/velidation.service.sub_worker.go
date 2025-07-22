package validation_service

import (
	"context"

	validation_adapters "github.com/root9464/Go_GamlerDefi/src/module/validation/adapters"
	validation_dto "github.com/root9464/Go_GamlerDefi/src/module/validation/dto"
	errors "github.com/root9464/Go_GamlerDefi/src/packages/lib/error"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *ValidationService) SubWorkerTransaction(ctx context.Context, transaction *validation_dto.WorkerTransactionDTO) (*validation_dto.WorkerTransactionDTO, bool, error) {
	s.logger.Info("start subworker transaction")
	s.logger.Infof("transaction: %+v", transaction)

	s.logger.Info("convert transaction id to bson.ObjectID")
	transactionID, err := bson.ObjectIDFromHex(transaction.ID)
	if err != nil {
		s.logger.Errorf("failed to convert transaction id to bson.ObjectID: %v", err)
		return transaction, false, err
	}

	s.logger.Info("precheckout transaction in db and update status to running")
	tr, err := s.validation_repository.PrecheckoutTransaction(ctx, transactionID)
	s.logger.Infof("transaction after precheckout in db: %+v", tr)
	transaction = validation_adapters.TransactionModelToDTOPoint(tr)
	if err != nil {
		s.logger.Errorf("failed to precheckout transaction: %v", err)
		return transaction, false, errors.NewError(400, err.Error())
	}

	s.logger.Info("precheckout transaction success, start next step...")
	return transaction, true, nil
}
