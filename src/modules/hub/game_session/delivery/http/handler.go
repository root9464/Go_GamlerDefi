package game_session_http

import (
	"github.com/gofiber/fiber/v2"
	game_session_hub "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/hub.go"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type Handler struct {
	hubManager *game_session_hub.Hub
	logger     *logger.Logger
}

func NewGameSessionHandler(hubManager *game_session_hub.Hub, logger *logger.Logger) *Handler {
	return &Handler{
		hubManager: hubManager,
		logger:     logger,
	}
}

// type CreateSessionRequest struct {
// 	GameName   string `json:"game_name"`
// 	PlayerID   string `json:"host_id"` // ID создателя сессии
// 	PlayerName string `json:"player_name"`
// }
//
// func (h *Handler) CreateSession(c *fiber.Ctx) error {
// 	dto := new(CreateSessionRequest)
// 	if err := c.BodyParser(dto); err != nil {
// 		return err
// 	}
//
// 	err := h.hubManager.CreateSession(c.Context(), dto.PlayerID, dto.PlayerName, dto.GameName)
// 	if err != nil {
// 		return err
// 	}
//
// 	return c.SendStatus(fiber.StatusOK)
// }

func (h *Handler) GetAllSessions(c *fiber.Ctx) error {
	sessions, err := h.hubManager.GetAllSessions(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(sessions)
}
