package game_session_ws

import "github.com/gofiber/fiber/v2"

func (h *GameSessionHandler) RegisterRoutes(router fiber.Router) {
	h.SetupSocketEventHandlers()
	router.Get("/ws/:session_id/:game_name/:user_id", h.GameSessionWSHandler)
}
