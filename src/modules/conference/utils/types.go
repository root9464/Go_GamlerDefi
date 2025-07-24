package conference_utils

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

type PeerConnection struct {
	PC   *webrtc.PeerConnection
	Conn *ThreadSafeWriter
}

type ThreadSafeWriter struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

type WebsocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func (t *ThreadSafeWriter) WriteJSON(v any) error {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return t.Conn.WriteJSON(v)
}
