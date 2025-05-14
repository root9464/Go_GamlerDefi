package referral_dto

// PaymentType defines types of referral payments
// @swagger:enum
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
	AuthorID int `json:"author_id,omitempty"`

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
	// ID of the author
	// required: true
	// minimum: 1
	// example: 12345
	AuthorID int `bson:"author_id" json:"author_id"`

	// ID of the referrer
	// required: true
	// minimum: 1
	// example: 12345
	ReferrerID int `bson:"referrer_id" json:"referrer_id"`

	// ID of the referral
	// required: true
	// minimum: 1
	// example: 12345
	ReferralID int `bson:"referral_id" json:"referral_id"`

	// Total amount of the payment
	// required: true
	// minimum: 0
	// example: 100.0
	TotalAmount float64 `bson:"total_amount" json:"total_amount"`

	// Number of tickets to process
	// required: true
	// minimum: 1
	// example: 100
	TicketCount int `bson:"ticket_count" json:"ticket_count"`

	// Levels of the payment
	// required: true
	// example: [{"level_number": 0, "rate": 0.2, "amount": 150, "address": "0QC3PUCoxBdLfOmO8xFQ84TGFPQUatxvvRsSAODKEvjbb4OS"}]
	Levels []LevelRequest `bson:"levels" json:"levels"`

	// Date of creation
	// required: true
	// example: 1715731200
	CreatedAt int64 `bson:"created_at" json:"created_at"`
}

// LevelRequest represents a level request
// @swagger:model LevelRequest
type LevelRequest struct {
	// Level number
	// required: true
	// minimum: 0
	// example: 0
	LevelNumber int `bson:"level_number"`

	// Rate of the level
	// required: true
	// minimum: 0
	// example: 0.2
	Rate float64 `bson:"rate"`

	// Amount of the level
	// required: true
	// minimum: 0
	// example: 150
	Amount float64 `bson:"amount"`

	// Address of the level
	// required: true
	// example: 0QC3PUCoxBdLfOmO8xFQ84TGFPQUatxvvRsSAODKEvjbb4OS
	Address string `bson:"address"`
}

// CellResponse represents a cell response
// @swagger:model CellResponse
type CellResponse struct {
	// Cell of the response
	// required: true
	// example: 0QC3PUCoxBdLfOmO8xFQ84TGFPQUatxvvRsSAODKEvjbb4OS
	Cell string `json:"cell"`
}
