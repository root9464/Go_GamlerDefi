package referral_service

import (
	"context"

	referral_repository "github.com/root9464/Go_GamlerDefi/modules/referral/repository"
	generated "github.com/root9464/Go_GamlerDefi/packages/generated/gql_generated"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
)

var _ IReferralService = &referralService{}

type IReferralService interface {
	CreateReferral(ctx context.Context, username string) (*generated.Referral, error)
}

type referralService struct {
	repo   referral_repository.IReferralRepository
	logger *logger.Logger
}

func NewReferralService(repo referral_repository.IReferralRepository, logger *logger.Logger) IReferralService {
	return &referralService{
		repo:   repo,
		logger: logger,
	}
}
