package game_config_handler

import (
	"github.com/gofiber/fiber/v2"
	game_config_dto "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/dto"
	game_config_service "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/service"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type GameConfigHandler struct {
	logger  *logger.Logger
	service *game_config_service.GameConfigService
}

func NewGameConfigHandler(logger *logger.Logger, service *game_config_service.GameConfigService) *GameConfigHandler {
	return &GameConfigHandler{logger: logger, service: service}
}

func (h *GameConfigHandler) CreateOrUpdate(ctx *fiber.Ctx) error {
	var dto game_config_dto.CreateGameConfigDTO
	if err := ctx.BodyParser(&dto); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.CreateOrUpdate(ctx.Context(), &dto); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "success"})
}
