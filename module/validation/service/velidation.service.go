package validation_service

import (
	validation_dto "github.com/root9464/Go_GamlerDefi/module/validation/dto"
	validation_repository "github.com/root9464/Go_GamlerDefi/module/validation/repository"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/tonkeeper/tonapi-go"
)

type IValidationService interface {
	RunnerTransaction(transaction *validation_dto.WorkerTransactionDTO) (*validation_dto.WorkerTransactionDTO, bool, error)
	SubWorkerTransaction(transaction *validation_dto.WorkerTransactionDTO) (*validation_dto.WorkerTransactionDTO, bool, error)
	WorkerTransaction(transaction *validation_dto.WorkerTransactionDTO) (*validation_dto.WorkerTransactionDTO, bool, error)
}

type ValidationService struct {
	logger  *logger.Logger
	ton_api *tonapi.Client

	validation_repository validation_repository.IValidationRepository
}

func NewValidationService(
	logger *logger.Logger, ton_api *tonapi.Client,
	validation_repository validation_repository.IValidationRepository,
) IValidationService {
	return &ValidationService{logger: logger, ton_api: ton_api, validation_repository: validation_repository}
}
