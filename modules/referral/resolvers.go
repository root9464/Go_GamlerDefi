package referral_module

var _ ReferralQueries = &ReferralResolver{}
var _ ReferralMutations = &ReferralResolver{}

type ReferralResolver struct{}

type ReferralQueries interface {
}

type ReferralMutations interface {
}
