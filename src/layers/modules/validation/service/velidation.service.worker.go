package validation_service

import (
	"context"
	"strings"
	"time"

	validation_adapters "github.com/root9464/Go_GamlerDefi/src/layers/modules/validation/adapters"
	validation_dto "github.com/root9464/Go_GamlerDefi/src/layers/modules/validation/dto"
	validation_model "github.com/root9464/Go_GamlerDefi/src/layers/modules/validation/model"
	errors "github.com/root9464/Go_GamlerDefi/src/packages/lib/error"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/address"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func IsTransactionValid(tx *tonapi.Trace) validation_dto.WorkerStatus {
	if tx.Transaction.ComputePhase.IsSet() &&
		!tx.Transaction.ComputePhase.Value.Skipped &&
		!tx.Transaction.ComputePhase.Value.Success.Value {
		return validation_dto.WorkerStatusFailed
	}

	if tx.Transaction.ActionPhase.IsSet() &&
		(!tx.Transaction.ActionPhase.Value.Success || tx.Transaction.ActionPhase.Value.ResultCode != 0) {
		return validation_dto.WorkerStatusWaiting
	}

	for _, child := range tx.Children {
		return IsTransactionValid(&child)
	}

	return validation_dto.WorkerStatusSuccess
}

func (s *ValidationService) IsAccountValid(transaction *validation_dto.WorkerTransactionDTO, trace *tonapi.Trace) bool {
	address, err := address.ParseAddr(transaction.TargetAddress)
	if err != nil {
		s.logger.Errorf("failed to parse user friendly address: %v", err)
		return false
	}

	s.logger.Infof("address: %+v", address.StringRaw())
	s.logger.Infof("destination: %+v", trace.Transaction.InMsg.Value.Destination.Value.Address)

	return strings.EqualFold(address.StringRaw(), trace.Transaction.InMsg.Value.Destination.Value.Address)
}

func (s *ValidationService) WorkerTransaction(ctx context.Context, transaction *validation_dto.WorkerTransactionDTO) (*validation_dto.WorkerTransactionDTO, bool, error) {
	s.logger.Info("start worker transaction")
	s.logger.Infof("transaction: %+v", transaction)
	transactionID, err := bson.ObjectIDFromHex(transaction.ID)
	if err != nil {
		s.logger.Errorf("failed to convert transaction id: %v", err)
		return nil, false, err
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	select {
	case <-ctx.Done():
		s.logger.Infof("transaction validation timed out")
		transaction, status, err := s.finalizeTransaction(ctx, transactionID, validation_dto.WorkerStatusFailed)
		if err != nil {
			s.logger.Errorf("failed to finalize transaction: %v", err)
			return transaction, false, err
		}
		s.logger.Infof("transaction validation timed out, status: %v", status)
		s.logger.Infof("transaction data: %+v", transaction)
		return transaction, status, errors.NewError(408, "transaction validation timed out")
	default:
		s.logger.Infof("get transaction trace by tx hash: %v", transaction.TxHash)
		txTrace, err := s.ton_api.GetTrace(ctx, tonapi.GetTraceParams{
			TraceID: transaction.TxHash,
		})

		if err != nil {
			s.logger.Errorf("failed to get transaction trace: %v", err)
			transaction, status, err := s.finalizeTransaction(ctx, transactionID, validation_dto.WorkerStatusFailed)
			if err != nil {
				s.logger.Errorf("failed to finalize transaction: %v", err)
				return transaction, false, err
			}
			s.logger.Infof("transaction finalize error, status: %v", status)
			s.logger.Infof("transaction data: %+v", transaction)
			return transaction, status, errors.NewError(502, "failed to get transaction trace")
		}

		s.logger.Infof("validate transaction: %v", transaction.TxHash)
		isAccountValid := s.IsAccountValid(transaction, txTrace)
		if !isAccountValid {
			s.logger.Errorf("account is not valid: %v", transaction.TargetAddress)
			transaction, status, err := s.finalizeTransaction(ctx, transactionID, validation_dto.WorkerStatusFailed)
			if err != nil {
				s.logger.Errorf("failed to finalize transaction: %v", err)
				return transaction, false, err
			}
			s.logger.Infof("transaction status updated to failed: %v", status)
			s.logger.Infof("transaction data: %+v", transaction)
			return transaction, status, errors.NewError(400, "account is not valid")
		}

		isValid := IsTransactionValid(txTrace)
		if isValid == validation_dto.WorkerStatusWaiting {
			s.logger.Warnf("transaction incomplete, waiting: %v", transaction.TxHash)
			transaction, status, err := s.finalizeTransaction(ctx, transactionID, validation_dto.WorkerStatusWaiting)
			if err != nil {
				s.logger.Errorf("failed to finalize transaction: %v", err)
				return transaction, false, err
			}
			s.logger.Infof("transaction status updated to waiting: %v", status)
			s.logger.Infof("transaction data: %+v", transaction)
			return transaction, status, nil
		}

		if isValid == validation_dto.WorkerStatusFailed {
			s.logger.Errorf("transaction failed: %v", transaction.TxHash)
			transaction, status, err := s.finalizeTransaction(ctx, transactionID, validation_dto.WorkerStatusFailed)
			if err != nil {
				s.logger.Errorf("failed to finalize transaction: %v", err)
				return transaction, false, err
			}
			s.logger.Infof("transaction status updated to failed: %v", status)
			s.logger.Infof("transaction data: %+v", transaction)
			return transaction, status, nil
		}

		s.logger.Infof("validate transaction success")
		transaction, status, err := s.finalizeTransaction(ctx, transactionID, validation_dto.WorkerStatusSuccess)
		if err != nil {
			s.logger.Errorf("failed to finalize transaction: %v", err)
			return transaction, false, err
		}
		s.logger.Infof("transaction status updated to success: %v", status)
		s.logger.Infof("transaction data: %+v", transaction)
		return transaction, status, nil
	}
}

func (s *ValidationService) finalizeTransaction(ctx context.Context, transactionID bson.ObjectID, status validation_dto.WorkerStatus) (*validation_dto.WorkerTransactionDTO, bool, error) {
	success := status == validation_dto.WorkerStatusSuccess

	tr, err := s.validation_repository.UpdateStatus(ctx, transactionID, validation_model.WorkerStatus(status))
	if err != nil {
		s.logger.Errorf("status update failed: %v", err)
		return nil, false, errors.NewError(400, "status update failed, transaction id: "+transactionID.Hex())
	}

	transaction := validation_adapters.TransactionModelToDTOPoint(tr)
	s.logger.Infof("finalize transaction, status: %v", tr.Status)
	s.logger.Infof("finalize transaction data: %+v", transaction)
	return transaction, success, nil
}
