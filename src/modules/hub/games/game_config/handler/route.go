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

	game.Post("/create", h.CreateGame)
	// game.Patch("/:game_name/settings", h.UpdateSettings)
	game.Put("/:game_name/settings", h.OverwriteSettings)
	game.Post("/assets", h.UploadAsset)

	game.Get("/assets/:name", h.GetAsset)
}
