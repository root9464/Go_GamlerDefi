package validation_service

import (
	"github.com/go-playground/validator/v10"
	"github.com/root9464/Go_GamlerDefi/config"
	validation_dto "github.com/root9464/Go_GamlerDefi/module/validation/dto"
	validation_repository "github.com/root9464/Go_GamlerDefi/module/validation/repository"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/tonkeeper/tonapi-go"
)

type IValidationService interface {
	RunnerTransaction(transaction *validation_dto.WorkerTransactionDTO) (*validation_dto.WorkerTransactionDTO, bool, error)
	SubWorkerTransaction(transaction *validation_dto.WorkerTransactionDTO) (bool, error)
	WorkerTransaction(trID string) (bool, error)
}

type ValidationService struct {
	logger    *logger.Logger
	config    *config.Config
	validator *validator.Validate
	ton_api   *tonapi.Client

	validation_repository validation_repository.IValidationRepository
}

func NewValidationService(
	logger *logger.Logger, config *config.Config, validator *validator.Validate, ton_api *tonapi.Client,
	validation_repository validation_repository.IValidationRepository,
) IValidationService {
	return &ValidationService{logger: logger, config: config, validator: validator, ton_api: ton_api, validation_repository: validation_repository}
}
