package conference_ws

import (
	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
)

func (h *WSHandler) RegisterRoutes(router fiber.Router) {
	ws := router.Group("/websocket")
	ws.Get("/", socketio.New(h.ConfirenceSocketHandler))
}
