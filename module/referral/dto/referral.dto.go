package referral_dto

type PaymentType string

const (
	PaymentAuthor   PaymentType = "become_author"
	PaymentReferred PaymentType = "join_game"
)

type ReferralProcessRequest struct {
	ReferrerID  int         `json:"referrer_id" validate:"required"`
	ReferredID  int         `json:"referred_id" validate:"required"`
	TicketCount int         `json:"ticket_count" validate:"required"`
	PaymentType PaymentType `json:"payment_type" validate:"required,oneof=become_author join_game"`
}
