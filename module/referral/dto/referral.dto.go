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

type PaymentOrder struct {
	AuthorID    int            `bson:"author_id" json:"author_id"`
	ReferrerID  int            `bson:"referrer_id" json:"referrer_id"`
	ReferralID  int            `bson:"referral_id" json:"referral_id"`
	TotalAmount float64        `bson:"total_amount" json:"total_amount"`
	TicketCount int            `bson:"ticket_count" json:"ticket_count"`
	CreatedAt   int64          `bson:"created_at" json:"created_at"`
	Levels      []LevelRequest `bson:"levels" json:"levels"`
}

type LevelRequest struct {
	LevelNumber int     `bson:"level_number"`
	Rate        float64 `bson:"rate"`
	Amount      float64 `bson:"amount"`
}
