package referral_dto

// PaymentType defines types of referral payments
// @swagger:enum
type PaymentType string

const (
	PaymentAuthor   PaymentType = "accrual_platform"
	PaymentReferred PaymentType = "leader_accrual"
)

// ReferralProcessRequest represents referral processing request
// @swagger:model ReferralProcessRequest
type ReferralProcessRequest struct {
	// ID of the referring user
	// required: true
	// minimum: 1
	// example: 12345
	ReferrerID int `json:"referrer_id" validate:"required"`

	// ID of the referred user
	// required: true
	// minimum: 1
	// example: 67890
	ReferredID int `json:"referred_id" validate:"required"`

	// Number of tickets to process
	// required: true
	// minimum: 1
	// example: 5
	TicketCount int `json:"ticket_count" validate:"required"`

	// Type of payment processing
	// required: true
	// enum: accrual_platform,leader_accrual
	// example: accrual_platform
	PaymentType PaymentType `json:"payment_type" validate:"required,oneof=accrual_platform leader_accrual"`
}
