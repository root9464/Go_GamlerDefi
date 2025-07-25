// package conference_utils
//
// import (
// 	"sync"
//
// 	"github.com/gofiber/contrib/websocket"
// 	"github.com/pion/webrtc/v4"
// )
//
// type PeerConnection struct {
// 	PC   *webrtc.PeerConnection
// 	Conn *ThreadSafeWriter
// }
//
// type ThreadSafeWriter struct {
// 	Conn *websocket.Conn
// 	Mu   sync.Mutex
// }
//
// type WebsocketMessage struct {
// 	Event string `json:"event"`
// 	Data  string `json:"data"`
// }
//
// func (t *ThreadSafeWriter) WriteJSON(v any) error {
// 	t.Mu.Lock()
// 	defer t.Mu.Unlock()
// 	return t.Conn.WriteJSON(v)
// }

package conference_utils

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/socketio"
	"github.com/pion/webrtc/v4"
)

type PeerConnection struct {
	PC   *webrtc.PeerConnection
	Conn *ThreadSafeWriter
}

type ThreadSafeWriter struct {
	Conn *socketio.Websocket
	Mu   sync.Mutex
}

type WebsocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func (t *ThreadSafeWriter) WriteJSON(v any) error {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	t.Conn.Emit(data, socketio.TextMessage)
	return nil
}
