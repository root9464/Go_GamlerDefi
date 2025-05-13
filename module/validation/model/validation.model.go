package validation_tr_model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WorkerStatus string

const (
	WorkerStatusPending  WorkerStatus = "pending"
	WorkerStatusWaiting  WorkerStatus = "waiting"
	WorkerStatusRunning  WorkerStatus = "running"
	WorkerStatusSuccess  WorkerStatus = "success"
	WorkerStatusFailed   WorkerStatus = "failed"
	WorkerStatusCanceled WorkerStatus = "canceled"
)

type WorkerTransaction struct {
	ID             string        `bson:"_id"`
	TxHash         string        `bson:"tx_hash"`
	TxBody         string        `bson:"tx_body"`
	TxQueryID      string        `bson:"tx_query_id"`
	TargetAddress  string        `bson:"target_address,omitempty"`
	PaymentOrderId bson.ObjectID `bson:"payment_order_id"`
	Status         WorkerStatus  `bson:"status"`
	CreatedAt      int64         `bson:"created_at"`
	UpdatedAt      int64         `bson:"updated_at"`
}
