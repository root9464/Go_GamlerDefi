package conference_utils

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/pion/webrtc/v4"
)

// PeerConnection - обертка для WebRTC соединения с дополнительной информацией.
type PeerConnection struct {
	PC     *webrtc.PeerConnection
	Writer *ThreadSafeWriter
	UserID string
	HubID  string
}

// ThreadSafeWriter обеспечивает потокобезопасную запись в WebSocket.
type ThreadSafeWriter struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

// WebsocketMessage - структура для обмена сообщениями через WebSocket.
type WebsocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// WriteJSON - потокобезопасный метод для отправки JSON.
func (t *ThreadSafeWriter) WriteJSON(v any) error {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return t.Conn.WriteJSON(v)
}
