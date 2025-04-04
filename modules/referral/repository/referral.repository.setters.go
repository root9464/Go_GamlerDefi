package referral_repository

import (
	"context"

	generated "github.com/root9464/Go_GamlerDefi/packages/generated/gql_generated"
)

const (
	referralCollection = "referrals"
)

func (r *referralRepository) CreateReferral(ctx context.Context, referral *generated.Referral) (*generated.Referral, error) {
	collection := r.mdb.Database("gamler_defi").Collection(referralCollection)
	result, err := collection.InsertOne(ctx, referral)
	if err != nil {
		r.logger.Errorf("Failed to insert referral into database: %v", err)
		return nil, err
	}

	r.logger.Infof("Successfully inserted referral: %v", result.InsertedID)

	return referral, nil
}
