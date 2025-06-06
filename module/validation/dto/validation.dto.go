package validation_dto

// WorkerStatus defines the status of the worker transaction
// @swagger:enum WorkerStatus
type WorkerStatus string

const (
	WorkerStatusPending WorkerStatus = "pending"
	WorkerStatusWaiting WorkerStatus = "waiting"
	WorkerStatusRunning WorkerStatus = "running"
	WorkerStatusSuccess WorkerStatus = "success"
	WorkerStatusFailed  WorkerStatus = "failed"
)

type WorkerTransactionDTO struct {
	// ID of the worker transaction
	// required: false
	// example: "682a67342a36c14af648479b"
	ID string `json:"id"`

	// Hash of the transaction
	// required: true
	// example: "105f7620bf78d534941ebcf97dda0dbe8e79c134a8ab346843787c71fe3308d5"
	TxHash string `json:"tx_hash" validate:"required"`

	// Query ID of the transaction
	// required: true
	// example: 1747000636
	TxQueryID uint64 `json:"tx_query_id" validate:"required"`

	// Target address of the transaction
	// required: true
	// example: "0QANsjLvOX2MERlT4oyv2bSPEVc9lunSPIs5a1kPthCXydUX"
	TargetAddress string `json:"target_address" validate:"required,len=48"`

	// Payment order ID
	// required: true
	// example: "6826ac79ff2f0eb00db5fa1d"
	PaymentOrderId string `json:"payment_order_id,omitempty"`

	// Status of the worker transaction
	// required: true
	// example: "pending"
	Status WorkerStatus `json:"status" validate:"required"`

	// Created at
	// required: false
	// example: 1715731200
	CreatedAt int64 `json:"created_at"`

	// Updated at
	// required: false
	// example: 1715731200
	UpdatedAt int64 `json:"updated_at"`
}

// WorkerTransactionResponse represents the response of the worker transaction
// @swagger:model WorkerTransactionResponse
type WorkerTransactionResponse struct {
	// Message of the response
	// example: "Transaction processed successfully"
	Message string `json:"message"`

	// Transaction of the response
	// example: "105f7620bf78d534941ebcf97dda0dbe8e79c134a8ab346843787c71fe3308d5"
	TxHash string `json:"tx_hash"`

	// Transaction ID
	// example: "682a67342a36c14af648479b"
	TxID string `json:"tx_id"`

	// Status of the worker transaction
	Status WorkerStatus `json:"status"`
}
