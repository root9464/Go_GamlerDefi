package conference_ws

import (
	"github.com/gofiber/fiber/v2"
)

func (h *WSHandler) RegisterRoutes(router fiber.Router) {
	ws := router.Group("/")
	ws.Get("/websocket", h.ConferenceWebsocketHandler)
}
