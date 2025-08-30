package game_session_dto

import "github.com/shopspring/decimal"

// test
type SessionInfo struct {
	ID          string `json:"id"`
	GameName    string `json:"gameName"`
	PlayerCount int    `json:"playerCount"`
	HostID      string `json:"hostId"`
}

type CreateSession struct {
	GameID           string          `json:"game_id"`
	StartTime        string          `json:"start_time"`
	Cost             decimal.Decimal `json:"cost"`
	Currency         string          `json:"currency"`
	TotalPlayers     int             `json:"total_players"`
	AvailablePlayers int             `json:"available_players"`
	Notes            string          `json:"notes"`
}
