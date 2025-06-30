package referral_dto

import "github.com/shopspring/decimal"

// PaymentType defines types of referral payments
// @swagger:enum PaymentType
type PaymentType string

const (
	PaymentPlatform PaymentType = "accrual_platform"
	PaymentLeader   PaymentType = "leader_accrual"
)

// ReferralProcessRequest represents referral processing request
// @swagger:model ReferralProcessRequest
type ReferralProcessRequest struct {
	// ID of the author
	// required: true
	// minimum: 1
	// example: 12345
	LeaderID int `json:"leader_id,omitempty"`

	// ID of the referring user
	// required: true
	// minimum: 1
	// example: 12345
	ReferrerID int `json:"referrer_id" validate:"required"`

	// ID of the referred user
	// required: true
	// minimum: 1
	// example: 67890
	ReferralID int `json:"referral_id" validate:"required"`

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

// PaymentOrder represents a payment order
// @swagger:model PaymentOrder
type PaymentOrder struct {
	// ID of the payment order
	// required: false
	// minimum: 1
	// example: 12345
	ID string `json:"id"`

	// ID of the author
	// required: true
	// minimum: 1
	// example: 12345
	LeaderID int `json:"leader_id"`

	// ID of the referrer
	// required: true
	// minimum: 1
	// example: 12345
	ReferrerID int `json:"referrer_id"`

	// ID of the referral
	// required: true
	// minimum: 1
	// example: 12345
	ReferralID int `json:"referral_id"`

	// Total amount of the payment
	// required: true
	// minimum: 0
	// example: 100.0
	TotalAmount decimal.Decimal `json:"total_amount"`

	// Number of tickets to process
	// required: true
	// minimum: 1
	// example: 100
	TicketCount int `json:"ticket_count"`

	// Levels of the payment
	// required: true
	// example: [{"level_number": 0, "rate": 0.2, "amount": 150, "address": "0QC3PUCoxBdLfOmO8xFQ84TGFPQUatxvvRsSAODKEvjbb4OS"}]
	Levels []LevelRequest `json:"levels"`

	// Date of creation
	// required: true
	// example: 1715731200
	CreatedAt int64 `json:"created_at"`
}

// LevelRequest represents a level request
// @swagger:model LevelRequest
type LevelRequest struct {
	// Level number
	// required: true
	// minimum: 0
	// example: 0
	LevelNumber int `json:"level_number"`

	// Rate of the level
	// required: true
	// minimum: 0
	// example: 0.2
	Rate decimal.Decimal `json:"rate"`

	// Amount of the level
	// required: true
	// minimum: 0
	// example: 150
	Amount decimal.Decimal `json:"amount"`

	// Address of the level
	// required: true
	// example: 0QC3PUCoxBdLfOmO8xFQ84TGFPQUatxvvRsSAODKEvjbb4OS
	Address string `json:"address"`
}

// CellResponse represents a cell response
// @swagger:model CellResponse
type CellResponse struct {
	// Cell of the response
	// required: true
	// example: 0QC3PUCoxBdLfOmO8xFQ84TGFPQUatxvvRsSAODKEvjbb4OS
	Cell string `json:"cell"`
}

// ValidateInvitationConditionsResponse represents a validate invitation conditions response
// @swagger:model ValidateInvitationConditionsResponse
type ValidateInvitationConditionsResponse struct {
	// Valid of the response
	// required: true
	// example: true
	Valid bool `json:"valid"`
}
