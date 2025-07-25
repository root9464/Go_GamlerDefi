package hub_entity

import (
	"errors"
	"slices"
	"time"
)

type HubStatus string
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
)

const (
	HubStatusNotStarted HubStatus = "not_started"
	HubStatusStarted    HubStatus = "started"
	HubStatusFinished   HubStatus = "finished"
)

type Hub struct {
	ID           string
	GameID       string   // Reference to the Game being played
	HostID       string   // ID of the host/moderator
	Participants []string // List of user IDs participating
	Status       HubStatus
	StartTime    time.Time // Scheduled start time
	EndTime      time.Time // When the game actually ended
	EntryFee     float64   // Cost to join
	Currency     Currency
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (h *Hub) IsActive() bool {
	now := time.Now()
	return h.Status == HubStatusStarted && !h.StartTime.After(now) && (h.EndTime.IsZero() || h.EndTime.After(now))
}

func (h *Hub) Start() error {
	if h.Status != HubStatusNotStarted {
		return errors.New("hub cannot be started: invalid status")
	}
	if time.Now().Before(h.StartTime) {
		return errors.New("hub cannot be started before scheduled time")
	}
	h.Status = HubStatusStarted
	h.UpdatedAt = time.Now()
	return nil
}

func (h *Hub) Finish() error {
	if h.Status != HubStatusStarted {
		return errors.New("hub cannot be finished: not started")
	}
	h.Status = HubStatusFinished
	h.EndTime = time.Now()
	h.UpdatedAt = time.Now()
	return nil
}

func (h *Hub) AddParticipant(userID string) error {
	if h.Status != HubStatusNotStarted {
		return errors.New("cannot add participant: hub already started")
	}

	if slices.Contains(h.Participants, userID) {
		return errors.New("user already in the hub")
	}

	h.Participants = append(h.Participants, userID)
	h.UpdatedAt = time.Now()
	return nil
}

func (h *Hub) RemoveParticipant(userID string) error {
	if h.Status != HubStatusNotStarted {
		return errors.New("cannot remove participant: hub already started")
	}
	for i, p := range h.Participants {
		if p == userID {
			h.Participants = append(h.Participants[:i], h.Participants[i+1:]...)
			h.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("user not found in participants")
}

func (h *Hub) HasParticipant(userID string) bool {
	return slices.Contains(h.Participants, userID)
}

func (h *Hub) SetCurrency(currency Currency) error {
	switch currency {
	case CurrencyUSD, CurrencyEUR, CurrencyGBP:
		h.Currency = currency
		h.UpdatedAt = time.Now()
		return nil
	default:
		return errors.New("unsupported currency")
	}
}
