package referral_dto

type PaymentTtpe string

const (
	PaymentAuthor   PaymentTtpe = "become_author"
	PaymentReferred PaymentTtpe = "join_game"
)

type ReferralProcessRequest struct {
	ReferrerID  int         `json:"referrer_id" validate:"required"`
	ReferredID  int         `json:"referred_id" validate:"required"`
	TicketCount int         `json:"ticket_count" validate:"required"`
	PaymentType PaymentTtpe `json:"payment_type" validate:"required,oneof=become_author join_game"`
}
