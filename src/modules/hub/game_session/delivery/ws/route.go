package game_session_ws

import "github.com/gofiber/fiber/v2"

func (h *GameSessionHandler) RegisterRoutes(router fiber.Router) {
	h.SetupSocketEventHandlers()
	session := router.Group("/session")
	session.Get("/ws/:session_id/:user_id", h.GameSessionWSHandler)
}
