package hub_entity

import "time"

type GameSession struct {
	ID           string
	HostID       string
	Participants []string
	StartTime    time.Time
	EndTime      time.Time
}
