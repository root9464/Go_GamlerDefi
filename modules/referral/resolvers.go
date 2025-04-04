package referral_resolvers

import (
	"context"

	referral_service "github.com/root9464/Go_GamlerDefi/modules/referral/service"
	generated "github.com/root9464/Go_GamlerDefi/packages/generated/gql_generated"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ ReferralQueries = &ReferralResolver{}
var _ ReferralMutations = &ReferralResolver{}

type ReferralResolver struct {
	mdb    *mongo.Client
	logger *logger.Logger

	referralService referral_service.IReferralService
}

type ReferralQueries interface {
	GetReferralByUsername(ctx context.Context, username string) (*generated.Referral, error)
	GetReferralByID(ctx context.Context, id int) (*generated.Referral, error)
}

type ReferralMutations interface {
	CreateReferral(ctx context.Context, username string) (*generated.Referral, error)
}

func (a *ReferralResolver) GetReferralByID(ctx context.Context, id int) (*generated.Referral, error) {
	panic("unimplemented")
}

func (a *ReferralResolver) GetReferralByUsername(ctx context.Context, username string) (*generated.Referral, error) {
	panic("unimplemented")
}

func (a *ReferralResolver) CreateReferral(ctx context.Context, username string) (*generated.Referral, error) {
	return a.referralService.CreateReferral(ctx, username)
}

func NewReferralResolver(mdb *mongo.Client, logger *logger.Logger) *ReferralResolver {
	return &ReferralResolver{mdb: mdb, logger: logger}
}
