package referral_model

import "time"

type PaymentOrder struct {
	AuthorID    int       `bson:"author_id"`
	ReferrerID  int       `bson:"referrer_id"`
	ReferralID  int       `bson:"referral_id"`
	TotalAmount float64   `bson:"total_amount"`
	TicketCount int       `bson:"ticket_count"`
	CreatedAt   time.Time `bson:"created_at"`
	Levels      []Level   `bson:"levels"`
}

type Level struct {
	LevelNumber int     `bson:"level_number"`
	Rate        float64 `bson:"rate"`
	Amount      float64 `bson:"amount"`
}
