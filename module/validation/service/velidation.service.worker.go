package validation_service

import (
	"context"
	"time"

	validation_dto "github.com/root9464/Go_GamlerDefi/module/validation/dto"
	validation_model "github.com/root9464/Go_GamlerDefi/module/validation/model"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"github.com/tonkeeper/tonapi-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *ValidationService) SubWorkerTransaction(transaction *validation_dto.WorkerTransactionDTO) (bool, error) {
	s.logger.Info("start worker transaction")
	s.logger.Infof("transaction: %+v", transaction)

	s.logger.Info("convert transaction id to bson.ObjectID")
	transactionID, err := bson.ObjectIDFromHex(transaction.ID)
	if err != nil {
		s.logger.Errorf("failed to convert transaction id to bson.ObjectID: %v", err)
		return false, err
	}

	s.logger.Info("precheckout transaction in db and update status to running")
	err = s.validation_repository.PrecheckoutTransaction(transactionID)
	if err != nil {
		s.logger.Errorf("failed to precheckout transaction: %v", err)
		return false, errors.NewError(400, err.Error())
	}

	s.logger.Info("precheckout transaction success, start next step...")
	return true, nil
}

const (
	maxRetries    = 12
	retryInterval = 5 * time.Second
	timeout       = 1 * time.Minute
)

func (s *ValidationService) WorkerTransaction(transaction *validation_dto.WorkerTransactionDTO) (bool, error) {
	s.logger.Info("start worker transaction")
	transactionID, err := bson.ObjectIDFromHex(transaction.ID)
	if err != nil {
		s.logger.Errorf("failed to convert transaction id: %v", err)
		return false, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for attempt := 0; attempt < maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			s.logger.Error("transaction validation timed out")
			return s.finalizeTransaction(transactionID, false)
		default:
		}

		if attempt == 0 {
			if err := s.validation_repository.UpdateStatus(transactionID,
				validation_model.WorkerStatusRunning); err != nil {
				return false, err
			}
		}

		txTrace, err := s.getTransactionTrace(transaction.TxHash)
		if err != nil {
			if attempt == 0 {
				err = s.validation_repository.UpdateStatus(transactionID, validation_model.WorkerStatus(validation_dto.WorkerStatusWaiting))
				return false, err
			}
			s.sleepWithContext(ctx, retryInterval)
			continue
		}

		isValid, err := s.ValidatorTransaction(transaction, txTrace)
		if err != nil {
			return false, err
		}

		if isValid {
			return s.finalizeTransaction(transactionID, true)
		}

		if attempt == 0 {
			err = s.validation_repository.UpdateStatus(transactionID, validation_model.WorkerStatus(validation_dto.WorkerStatusWaiting))
			return false, err
		}
		s.sleepWithContext(ctx, retryInterval)
	}

	return s.finalizeTransaction(transactionID, false)
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

func (s *ValidationService) finalizeTransaction(transactionID bson.ObjectID, success bool) (bool, error) {
	status := validation_dto.WorkerStatusFailed
	if success {
		status = validation_dto.WorkerStatusSuccess
	}

	if err := s.validation_repository.UpdateStatus(transactionID, validation_model.WorkerStatus(status)); err != nil {
		s.logger.Errorf("status update failed: %v", err)
		return false, err
	}

	return success, nil
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
