package referral_model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PaymentOrder struct {
	ID          bson.ObjectID   `bson:"_id"`
	LeaderID    int             `bson:"leader_id"`
	ReferrerID  int             `bson:"referrer_id"`
	ReferralID  int             `bson:"referral_id"`
	TotalAmount bson.Decimal128 `bson:"total_amount"`
	TicketCount int             `bson:"ticket_count"`
	CreatedAt   int64           `bson:"created_at"`
	TrHash      string          `bson:"tr_hash,omitempty"`
	Levels      []Level         `bson:"levels"`
}

type Level struct {
	LevelNumber int             `bson:"level_number"`
	Rate        bson.Decimal128 `bson:"rate"`
	Amount      bson.Decimal128 `bson:"amount"`
	Address     string          `bson:"address"`
}
