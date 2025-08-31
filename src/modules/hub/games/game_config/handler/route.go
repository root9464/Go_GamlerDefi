package game_config_handler

import "github.com/gofiber/fiber/v2"

// func (h *Handler) RegisterRoutes(router fiber.Router) {
// 	game := router.Group("/game")
//
// 	game.Get("/sessions", h.GetAllSessions)
// 	// game.Post("/session", h.CreateSession)
// }

func (h *GameConfigHandler) RegisterRoutes(router fiber.Router) {
	game := router.Group("/game")
	game.Post("/config", h.CreateOrUpdate)
}
