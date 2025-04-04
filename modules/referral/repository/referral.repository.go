package referral_repository

import (
	"context"

	generated "github.com/root9464/Go_GamlerDefi/packages/generated/gql_generated"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ IReferralRepository = &referralRepository{}

type IReferralRepository interface {
	CreateReferral(ctx context.Context, referral *generated.Referral) (*generated.Referral, error)
}

type referralRepository struct {
	mdb    *mongo.Client
	logger *logger.Logger
}

func NewReferralRepository(mdb *mongo.Client, logger *logger.Logger) IReferralRepository {
	return &referralRepository{
		mdb:    mdb,
		logger: logger,
	}
}
