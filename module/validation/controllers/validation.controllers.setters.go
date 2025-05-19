package validation_controllers

import (
	"github.com/gofiber/fiber/v2"
	validation_dto "github.com/root9464/Go_GamlerDefi/module/validation/dto"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
)

func (c *ValidationController) ValidatorTransaction(ctx *fiber.Ctx) error {
	transaction := new(validation_dto.WorkerTransactionDTO)
	if err := ctx.BodyParser(transaction); err != nil {
		c.logger.Errorf("failed to parse transaction: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	c.logger.Info("validate dto")
	if err := c.validator.Struct(transaction); err != nil {
		c.logger.Errorf("failed to validate transaction dto: %v", err)
		return errors.NewError(400, "invalid request body")
	}
	c.logger.Info("validate dto success")

	transaction, runnerStatus, err := c.validation_service.RunnerTransaction(transaction)
	if err != nil {
		c.logger.Errorf("runner failed: %v", err)
		return err
	}

	if !runnerStatus {
		c.logger.Errorf("failed runner transaction: %v", transaction.TxHash)
		return ctx.Status(fiber.StatusConflict).JSON(validation_dto.WorkerTransactionResponse{
			Message: "Failed transaction processing",
			TxHash:  transaction.TxHash,
			TxID:    transaction.ID,
			Status:  transaction.Status,
		})
	}

	transaction, subWorkerStatus, err := c.validation_service.SubWorkerTransaction(transaction)
	if err != nil {
		c.logger.Errorf("failed subworker transaction: %v", err)
		return err
	}

	if !subWorkerStatus {
		c.logger.Errorf("failed subworker transaction: %v", transaction.TxHash)
		return ctx.Status(fiber.StatusConflict).JSON(validation_dto.WorkerTransactionResponse{
			Message: "Failed subworker transaction",
			TxHash:  transaction.TxHash,
			TxID:    transaction.ID,
			Status:  transaction.Status,
		})
	}

	transaction, workerStatus, err := c.validation_service.WorkerTransaction(transaction)
	if err != nil {
		c.logger.Errorf("failed worker transaction: %v", err)
		return err
	}

	if !workerStatus {
		c.logger.Errorf("failed worker transaction: %v", transaction.TxHash)
		return ctx.Status(fiber.StatusConflict).JSON(validation_dto.WorkerTransactionResponse{
			Message: "Failed worker transaction",
			TxHash:  transaction.TxHash,
			TxID:    transaction.ID,
			Status:  transaction.Status,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(validation_dto.WorkerTransactionResponse{
		Message: "Transaction processed successfully",
		TxHash:  transaction.TxHash,
		TxID:    transaction.ID,
		Status:  transaction.Status,
	})
}
