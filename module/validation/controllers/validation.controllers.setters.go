package validation_controllers

import (
	"github.com/gofiber/fiber/v2"
	validation_dto "github.com/root9464/Go_GamlerDefi/module/validation/dto"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
)

// @Summary Validate transaction
// @Description Validate transaction
// @Tags Validation
// @Accept json
// @Produce json
// @Param transaction body validation_dto.WorkerTransactionDTO true "Transaction"
// @Success 200 {object} validation_dto.WorkerTransactionResponse "Transaction processed successfully"
// @Failure 400 {object} errors.MapError "Invalid request body"
// @Failure 409 {object} validation_dto.WorkerTransactionResponse "Failed transaction processing"
// @Router /api/validation/validate [post]
func (c *ValidationController) ValidatorTransaction(ctx *fiber.Ctx) error {
	transaction := new(validation_dto.WorkerTransactionDTO)
	if err := ctx.BodyParser(transaction); err != nil {
		c.logger.Errorf("failed to parse transaction: %v", err)
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	c.logger.Info("validate dto")
	if err := c.validator.Struct(transaction); err != nil {
		c.logger.Errorf("failed to validate transaction dto: %v", err)
		return errors.NewError(400, "invalid request body")
	}
	c.logger.Info("validate dto success")

	transaction, runnerStatus, err := c.validation_service.RunnerTransaction(ctx.Context(), transaction)
	if err != nil {
		c.logger.Errorf("runner failed: %v", err)
		return ctx.Status(errors.GetCode(err)).JSON(validation_dto.WorkerTransactionResponse{
			Message: err.Error(),
			TxHash:  transaction.TxHash,
			TxID:    transaction.ID,
			Status:  transaction.Status,
		})
	}

	if !runnerStatus {
		c.logger.Errorf("failed runner transaction: %v", transaction.TxHash)
		return ctx.Status(409).JSON(validation_dto.WorkerTransactionResponse{
			Message: "Failed transaction processing",
			TxHash:  transaction.TxHash,
			TxID:    transaction.ID,
			Status:  transaction.Status,
		})
	}

	c.logger.Infof("transaction data after runner: %+v", transaction)

	transaction, subWorkerStatus, err := c.validation_service.SubWorkerTransaction(ctx.Context(), transaction)
	if err != nil {
		c.logger.Errorf("failed subworker transaction: %v", err)
		return ctx.Status(errors.GetCode(err)).JSON(validation_dto.WorkerTransactionResponse{
			Message: err.Error(),
			TxHash:  transaction.TxHash,
			TxID:    transaction.ID,
			Status:  transaction.Status,
		})
	}
	c.logger.Infof("transaction data after subworker: %+v", transaction)

	if !subWorkerStatus {
		c.logger.Errorf("failed subworker transaction: %v", transaction.TxHash)
		return ctx.Status(409).JSON(validation_dto.WorkerTransactionResponse{
			Message: "Failed subworker transaction",
			TxHash:  transaction.TxHash,
			TxID:    transaction.ID,
			Status:  transaction.Status,
		})
	}

	c.logger.Infof("transaction data before worker: %+v", transaction)
	transaction, workerStatus, err := c.validation_service.WorkerTransaction(ctx.Context(), transaction)
	if err != nil {
		c.logger.Errorf("failed worker transaction: %v", err)
		return ctx.Status(errors.GetCode(err)).JSON(validation_dto.WorkerTransactionResponse{
			Message: err.Error(),
			TxHash:  transaction.TxHash,
			TxID:    transaction.ID,
			Status:  transaction.Status,
		})
	}

	if !workerStatus {
		c.logger.Errorf("failed worker transaction: %v", transaction.TxHash)
		return ctx.Status(errors.GetCode(err)).JSON(validation_dto.WorkerTransactionResponse{
			Message: "Failed worker transaction",
			TxHash:  transaction.TxHash,
			TxID:    transaction.ID,
			Status:  transaction.Status,
		})
	}

	return ctx.Status(200).JSON(validation_dto.WorkerTransactionResponse{
		Message: "Transaction processed successfully",
		TxHash:  transaction.TxHash,
		TxID:    transaction.ID,
		Status:  transaction.Status,
	})
}
