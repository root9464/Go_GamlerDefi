package hub_entity

import "time"

type TicketStatus string

const (
	TicketStatusActive TicketStatus = "Активен"
	TicketStatusUsed   TicketStatus = "Потрачен"
	TicketStatusBooked TicketStatus = "Забронировано"
)

type Ticket struct {
	ID             string
	TicketStatus   TicketStatus
	DateOfIssue    time.Time
	ExpirationDate time.Time
	GameID         string
}
