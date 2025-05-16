package validation_dto

type WorkerStatus string

const (
	WorkerStatusPending WorkerStatus = "pending"
	WorkerStatusWaiting WorkerStatus = "waiting"
	WorkerStatusRunning WorkerStatus = "running"
	WorkerStatusSuccess WorkerStatus = "success"
	WorkerStatusFailed  WorkerStatus = "failed"
)

type WorkerTransactionDTO struct {
	ID                 string       `json:"id" validate:"required"`
	TxHash             string       `json:"tx_hash" validate:"required"`
	TxQueryID          uint64       `json:"tx_query_id" validate:"required"`
	TargetJettonSymbol string       `json:"target_jetton_symbol" validate:"required"`
	TargetJettonMaster string       `json:"target_jetton_master" validate:"required,len=48"`
	TargetAddress      string       `json:"target_address" validate:"required,len=48"`
	PaymentOrderId     string       `json:"payment_order_id" validate:"required"`
	Status             WorkerStatus `json:"status" validate:"required"`
	CreatedAt          int64        `json:"created_at"`
	UpdatedAt          int64        `json:"updated_at"`
}
