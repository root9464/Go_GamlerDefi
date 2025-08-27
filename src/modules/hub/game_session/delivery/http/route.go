package game_session_http

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	game := router.Group("/game")

	game.Get("/sessions", h.GetAllSessions)
	game.Post("/session", h.CreateSession)
}
