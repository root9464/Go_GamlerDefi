package referral_dto

type PaymentType string

const (
	PaymentAuthor   PaymentType = "accrual_platform"
	PaymentReferred PaymentType = "leader_accrual"
)

type ReferralProcessRequest struct {
	ReferrerID  int         `json:"referrer_id" validate:"required"`
	ReferredID  int         `json:"referred_id" validate:"required"`
	TicketCount int         `json:"ticket_count" validate:"required"`
	PaymentType PaymentType `json:"payment_type" validate:"required,oneof=accrual_platform leader_accrual"`
}
