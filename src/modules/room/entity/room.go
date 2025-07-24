package room_entity

import "time"

type RoomStatus string

const (
	RoomStatusNotStarted RoomStatus = "not_started"
	RoomStatusStarted    RoomStatus = "started"
	RoomStatusFinished   RoomStatus = "finished"
)

type Room struct {
	ID           string
	GameID       string   // Reference to the Game being played
	HostID       string   // ID of the host/moderator
	Participants []string // List of user IDs participating
	Status       RoomStatus
	StartTime    time.Time // Scheduled start time
	EndTime      time.Time // When the game actually ended
	MaxSeats     int       // Maximum participants (could be derived from Game.MaxPlayers)
	EntryFee     float64   // Cost to join
	// VideoConfURL string    // Video conference link
	CreatedAt time.Time
	UpdatedAt time.Time
}
