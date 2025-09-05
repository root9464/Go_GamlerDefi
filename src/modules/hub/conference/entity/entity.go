package conference_entity

import (
	"sync"
	"time"

	"github.com/gofiber/contrib/socketio"
	"github.com/pion/webrtc/v4"
)

type Connection struct {
	Pc             *webrtc.PeerConnection
	Kws            *socketio.Websocket
	Lock           sync.Mutex
	Tracks         map[string]*webrtc.TrackLocalStaticRTP
	RoomID         string
	User           *User
	Closed         chan struct{}
	LastSignal     time.Time
	SignalDebounce time.Duration
}

type Room struct {
	Lock        sync.RWMutex
	Connections []*Connection
	TrackLocals map[string]*webrtc.TrackLocalStaticRTP
	TrackCount  int64
}

type User struct {
	StreamID string
	TrackID  string
	UserID   string
	Username string
}

func (r *Room) GetRoom(rooms map[string]*Room, roomID string) *Room {
	room, exists := rooms[roomID]
	if !exists {
		return nil
	}
	return room
}
