package validation_service

import (
	"context"
	"time"

	validation_adapters "github.com/root9464/Go_GamlerDefi/module/validation/adapters"
	validation_dto "github.com/root9464/Go_GamlerDefi/module/validation/dto"
	validation_model "github.com/root9464/Go_GamlerDefi/module/validation/model"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"github.com/tonkeeper/tonapi-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *ValidationService) SubWorkerTransaction(transaction *validation_dto.WorkerTransactionDTO) (*validation_dto.WorkerTransactionDTO, bool, error) {
	s.logger.Info("start worker transaction")
	s.logger.Infof("transaction: %+v", transaction)

	s.logger.Info("convert transaction id to bson.ObjectID")
	transactionID, err := bson.ObjectIDFromHex(transaction.ID)
	if err != nil {
		s.logger.Errorf("failed to convert transaction id to bson.ObjectID: %v", err)
		return nil, false, err
	}

	s.logger.Info("precheckout transaction in db and update status to running")
	tr, err := s.validation_repository.PrecheckoutTransaction(transactionID)
	transactionDTO := validation_adapters.TransactionModelToDTOPoint(tr)
	if err != nil {
		s.logger.Errorf("failed to precheckout transaction: %v", err)
		return nil, false, errors.NewError(400, err.Error())
	}

	s.logger.Info("precheckout transaction success, start next step...")
	return transactionDTO, true, nil
}

const (
	maxRetries    = 12
	retryInterval = 5 * time.Second
	timeout       = 1 * time.Minute
)

func (s *ValidationService) WorkerTransaction(transaction *validation_dto.WorkerTransactionDTO) (*validation_dto.WorkerTransactionDTO, bool, error) {
	s.logger.Info("start worker transaction")
	transactionID, err := bson.ObjectIDFromHex(transaction.ID)
	if err != nil {
		s.logger.Errorf("failed to convert transaction id: %v", err)
		return nil, false, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.logger.Infof("start worker transaction with %d attempts", maxRetries)
	for attempt := range maxRetries {
		s.logger.Infof("attempt: %d", attempt)
		select {
		case <-ctx.Done():
			s.logger.Infof("done transaction context")
			transaction, status, err := s.finalizeTransaction(transactionID, false)
			s.logger.Infof("transaction validation timed out, status: %v", status)
			s.logger.Infof("transaction data: %+v", transaction)
			return transaction, status, err
		default:
		}

		s.logger.Infof("get transaction trace by tx hash: %v", transaction.TxHash)
		txTrace, err := s.getTransactionTrace(transaction.TxHash)
		if err != nil {
			if attempt == 0 {
				tr, err := s.validation_repository.UpdateStatus(transactionID, validation_model.WorkerStatus(validation_dto.WorkerStatusWaiting))
				s.logger.Infof("transaction status updated to waiting: %v", tr.Status)
				transaction = validation_adapters.TransactionModelToDTOPoint(tr)
				s.logger.Infof("transaction data: %+v", transaction)
				return transaction, false, err
			}

			s.logger.Infof("sleep with context: %v", retryInterval)
			s.sleepWithContext(ctx, retryInterval)
			continue
		}

		s.logger.Infof("validate transaction: %v", transaction.TxHash)
		updatedTransaction, isValid, err := s.ValidatorTransaction(transaction, txTrace)
		transaction = updatedTransaction
		if err != nil {
			s.logger.Errorf("failed to validate transaction: %v", err)
			tr, err := s.validation_repository.UpdateStatus(transactionID, validation_model.WorkerStatus(validation_dto.WorkerStatusFailed))
			transaction = validation_adapters.TransactionModelToDTOPoint(tr)
			s.logger.Infof("transaction data: %+v", transaction)
			return transaction, false, err
		}

		s.logger.Infof("transaction is valid: %v", isValid)
		if isValid {
			s.logger.Infof("transaction status updated to success: %v", validation_dto.WorkerStatusSuccess)
			tr, err := s.validation_repository.UpdateStatus(transactionID, validation_model.WorkerStatus(validation_dto.WorkerStatusSuccess))
			transaction = validation_adapters.TransactionModelToDTOPoint(tr)
			s.logger.Infof("transaction data: %+v", transaction)
			return transaction, true, err
		}

		if attempt == 0 {
			s.logger.Infof("transaction status updated to waiting: %v", validation_dto.WorkerStatusWaiting)
			tr, err := s.validation_repository.UpdateStatus(transactionID, validation_model.WorkerStatus(validation_dto.WorkerStatusWaiting))
			transaction = validation_adapters.TransactionModelToDTOPoint(tr)
			s.logger.Infof("transaction data: %+v", transaction)
			return transaction, false, err
		}

		s.logger.Infof("sleep with context: %v", retryInterval)
		s.sleepWithContext(ctx, retryInterval)
	}

	transaction, status, err := s.finalizeTransaction(transactionID, false)
	s.logger.Infof("transaction status updated to failed: %v", validation_dto.WorkerStatusFailed)
	s.logger.Infof("transaction data: %+v", transaction)
	return transaction, status, err
}

func (s *ValidationService) getTransactionTrace(txHash string) (*tonapi.Trace, error) {
	txTrace, err := s.ton_api.GetTrace(context.Background(), tonapi.GetTraceParams{
		TraceID: txHash,
	})
	if err != nil {
		s.logger.Warnf("failed to get transaction trace: %v", err)
		return nil, err
	}
	return txTrace, nil
}

func (s *ValidationService) finalizeTransaction(transactionID bson.ObjectID, success bool) (*validation_dto.WorkerTransactionDTO, bool, error) {
	status := validation_dto.WorkerStatusFailed
	if success {
		status = validation_dto.WorkerStatusSuccess
	}

	tr, err := s.validation_repository.UpdateStatus(transactionID, validation_model.WorkerStatus(status))
	if err != nil {
		s.logger.Errorf("status update failed: %v", err)
		return nil, false, err
	}

	transaction := validation_adapters.TransactionModelToDTOPoint(tr)
	return transaction, success, nil
}

func (s *ValidationService) sleepWithContext(ctx context.Context, duration time.Duration) {
	timer := time.NewTimer(duration)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
	}
}
