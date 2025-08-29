package hub_entity

import "time"

type GameStatus string

const (
	StatusNotStarted GameStatus = "not_started"
	StatusStarted    GameStatus = "started"
	StatusFinished   GameStatus = "finished"
)

type Game struct {
	ID              string
	Name            string
	MaxPlayers      int
	DurationMinutes int
	Description     string
	Status          GameStatus
}

func (g *Game) IsValidStatus(status GameStatus) bool {
	switch status {
	case StatusNotStarted, StatusStarted, StatusFinished:
		return true
	default:
		return false
	}
}

func (g *Game) CanAddPlayer(currentPlayers int) bool {
	return currentPlayers < g.MaxPlayers
}

func (g *Game) IsFull(currentPlayers int) bool {
	return currentPlayers >= g.MaxPlayers
}

func (g *Game) EstimatedEndTime(startTime time.Time) time.Time {
	return startTime.Add(time.Duration(g.DurationMinutes) * time.Minute)
}

func (g *Game) TimeRemaining(startTime time.Time) time.Duration {
	if g.Status != StatusStarted {
		return 0
	}
	remaining := g.EstimatedEndTime(startTime).Sub(time.Now())
	if remaining < 0 {
		return 0
	}
	return remaining
}
