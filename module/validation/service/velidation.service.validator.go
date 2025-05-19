package validation_service

import (
	"strings"

	validation_dto "github.com/root9464/Go_GamlerDefi/module/validation/dto"
	validation_model "github.com/root9464/Go_GamlerDefi/module/validation/model"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"github.com/samber/lo"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/address"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func IsTransactionValid(tx *tonapi.Trace) bool {
	if tx.Transaction.ComputePhase.IsSet() &&
		!tx.Transaction.ComputePhase.Value.Skipped &&
		!tx.Transaction.ComputePhase.Value.Success.Value {
		return false
	}

	if tx.Transaction.ActionPhase.IsSet() &&
		(!tx.Transaction.ActionPhase.Value.Success || tx.Transaction.ActionPhase.Value.ResultCode != 0) {
		return false
	}

	return lo.EveryBy(tx.Children, func(child tonapi.Trace) bool {
		return IsTransactionValid(&child)
	})

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

func (s *ValidationService) ValidatorTransaction(transaction *validation_dto.WorkerTransactionDTO, tx *tonapi.Trace) (bool, error) {
	s.logger.Infof("starting validator transaction: %+v", transaction.TxHash)
	txHash := tx.Transaction.InMsg.Value.Hash
	s.logger.Infof("transaction in blockchain hash: %+v", txHash)
	s.logger.Infof("transaction in dto hash: %+v", transaction.TxHash)
	transactionID, err := bson.ObjectIDFromHex(transaction.ID)
	if err != nil {
		s.logger.Errorf("failed to convert transaction id to bson.ObjectID: %v", err)
		return false, err
	}

	if transaction.TxHash != txHash {
		s.logger.Errorf("transaction hash is not valid: %+v", txHash)
		return false, errors.NewError(400, "transaction hash is not valid")
	}

	isValid := IsTransactionValid(tx)
	if !isValid {
		s.logger.Warnf("transaction incomplete, waiting: %v", transaction.TxHash)
		_, err := s.validation_repository.UpdateStatus(transactionID, validation_model.WorkerStatus(validation_dto.WorkerStatusWaiting))
		if err != nil {
			s.logger.Errorf("failed to update status: %v", err)
			return false, err
		}
		return false, errors.NewError(409, "transaction processing not completed")
	}

	isAccountValid := s.IsAccountValid(transaction, tx)
	if !isAccountValid {
		s.logger.Errorf("account is not valid: %v", transaction.TargetAddress)
		_, err = s.validation_repository.UpdateStatus(transactionID, validation_model.WorkerStatus(validation_dto.WorkerStatusFailed))
		if err != nil {
			s.logger.Errorf("failed to update status: %v", err)
			return false, err
		}
		return false, errors.NewError(400, "account is not valid")
	}

	return true, nil
}
