package game_session_dto

// test
type SessionInfo struct {
	ID          string `json:"id"`
	GameName    string `json:"gameName"`
	PlayerCount int    `json:"playerCount"`
	HostID      string `json:"hostId"`
}
