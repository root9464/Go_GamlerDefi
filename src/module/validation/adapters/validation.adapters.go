package validation_adapters

import (
	validation_dto "github.com/root9464/Go_GamlerDefi/src/module/validation/dto"
	validation_model "github.com/root9464/Go_GamlerDefi/src/module/validation/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TransactionDTOToModel(transactionDTO validation_dto.WorkerTransactionDTO) (validation_model.WorkerTransaction, error) {
	transactionID, err := bson.ObjectIDFromHex(transactionDTO.ID)
	if err != nil {
		return validation_model.WorkerTransaction{}, err
	}

	var paymentOrderID bson.ObjectID
	if transactionDTO.PaymentOrderId != "" {
		paymentOrderID, err = bson.ObjectIDFromHex(transactionDTO.PaymentOrderId)
		if err != nil {
			return validation_model.WorkerTransaction{}, err
		}
	}

	return validation_model.WorkerTransaction{
		ID:             transactionID,
		TxHash:         transactionDTO.TxHash,
		TxQueryID:      transactionDTO.TxQueryID,
		TargetAddress:  transactionDTO.TargetAddress,
		PaymentOrderId: paymentOrderID,
		Status:         validation_model.WorkerStatus(transactionDTO.Status),
		CreatedAt:      transactionDTO.CreatedAt,
		UpdatedAt:      transactionDTO.UpdatedAt,
	}, nil
}

func TransactionModelToDTOPoint(transactionModel validation_model.WorkerTransaction) *validation_dto.WorkerTransactionDTO {
	return &validation_dto.WorkerTransactionDTO{
		ID:             transactionModel.ID.Hex(),
		TxHash:         transactionModel.TxHash,
		TxQueryID:      transactionModel.TxQueryID,
		TargetAddress:  transactionModel.TargetAddress,
		PaymentOrderId: transactionModel.PaymentOrderId.Hex(),
		Status:         validation_dto.WorkerStatus(transactionModel.Status),
		CreatedAt:      transactionModel.CreatedAt,
		UpdatedAt:      transactionModel.UpdatedAt,
	}
}
