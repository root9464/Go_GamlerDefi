package conference_ws

import "github.com/gofiber/contrib/socketio"

func (h *WSHandler) socketErr(ep *socketio.EventPayload, err error) {
	errByte := []byte(err.Error())
	ep.Kws.Emit(errByte)
}

func (h *WSHandler) WS() func(*socketio.Websocket) { return nil }
