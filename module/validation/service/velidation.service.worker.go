package validation_service

import (
	"context"
	"time"

	validation_dto "github.com/root9464/Go_GamlerDefi/module/validation/dto"
	validation_model "github.com/root9464/Go_GamlerDefi/module/validation/model"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"github.com/samber/lo"
	"github.com/tonkeeper/tonapi-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *ValidationService) SubWorkerTransaction(transaction *validation_dto.WorkerTransactionDTO) (bool, error) {
	s.logger.Info("start worker transaction")
	s.logger.Infof("transaction: %+v", transaction)

	s.logger.Info("validate dto")
	if err := s.validator.Struct(transaction); err != nil {
		s.logger.Errorf("failed to validate transaction dto: %v", err)
		return false, err
	}
	s.logger.Info("validate dto success")

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

func (s *ValidationService) WorkerTransaction(trID string) (bool, error) {
	s.logger.Info("start worker transaction")
	s.logger.Infof("transaction: %+v", trID)

	transactionID, err := bson.ObjectIDFromHex(trID)
	if err != nil {
		s.logger.Errorf("failed to convert transaction id to bson.ObjectID: %v", err)
		return false, err
	}

	s.logger.Info("get transaction from db")
	transaction, err := s.validation_repository.GetTransactionObserver(transactionID)
	if err != nil {
		s.logger.Errorf("failed to get transaction from database: %v", err)
		return false, err
	}

	s.logger.Infof("check transaction in blockchain")
	eventTx, err := s.ton_api.GetAccountEvent(context.Background(), tonapi.GetAccountEventParams{
		AccountID: transaction.TargetAddress,
		EventID:   transaction.TxHash,
	})

	if err != nil {
		s.logger.Errorf("failed to get transaction, try again later from blockchain, txHash: %s, error: %v", transaction.TxHash, err)
		err = s.validation_repository.UpdateStatus(transactionID, validation_model.WorkerStatus(validation_dto.WorkerStatusWaiting))
		if err != nil {
			s.logger.Errorf("failed to update status: %v", err)
			return false, err
		}
		return false, err
	}

	isAddressValid := eventTx.Account.Address == transaction.TargetAddress
	if !isAddressValid {
		s.logger.Errorf("transaction address is not valid: %s", eventTx.Account.Address)
		return false, errors.NewError(400, "transaction address is not valid")
	}

	startTime := time.Unix(int64(transaction.TxQueryID), 0)
	endTime := time.Unix(eventTx.Timestamp, 0)
	timeDiff := endTime.Sub(startTime)

	if timeDiff > time.Hour || timeDiff < -time.Hour {
		s.logger.Errorf("transaction time difference too large: start %s, end %s, diff %s", startTime, endTime, timeDiff)
		return false, errors.NewError(400, "transaction time difference exceeds 1 hour")
	}

	isOptsValid := lo.ContainsBy(eventTx.Actions, func(action tonapi.Action) bool {
		isSymbolValid := action.JettonTransfer.Value.Jetton.Symbol == transaction.TargetJettonSymbol
		isMasterValid := action.JettonTransfer.Value.Jetton.Address == transaction.TargetJettonMaster
		isStatusValid := action.Status == "ok"
		return isSymbolValid && isMasterValid && isStatusValid
	})

	if !isOptsValid {
		s.logger.Errorf("transaction opts is not valid: %s", eventTx.Actions[0].JettonTransfer.Value.Jetton.Symbol)
		err = s.validation_repository.UpdateStatus(transactionID, validation_model.WorkerStatus(validation_dto.WorkerStatusFailed))
		if err != nil {
			s.logger.Errorf("failed to update status: %v", err)
			return false, err
		}
		return false, errors.NewError(400, "transaction opts is not valid")
	}

	return true, nil
}
