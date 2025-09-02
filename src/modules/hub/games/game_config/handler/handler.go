package game_config_handler

import (
	"io"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/gofiber/fiber/v2"
	game_config_dto "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/dto"
	game_config_service "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/service"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	file_service "github.com/root9464/Go_GamlerDefi/src/submodules/file/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameConfigHandler struct {
	logger      *logger.Logger
	service     *game_config_service.GameConfigService
	fileService *file_service.FileService
}

func NewGameConfigHandler(logger *logger.Logger, service *game_config_service.GameConfigService, fileService *file_service.FileService) *GameConfigHandler {
	return &GameConfigHandler{logger: logger, service: service, fileService: fileService}
)

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

// ==============================================================
func (h *GameConfigHandler) CreateGame(c *fiber.Ctx) error {
	var dto game_config_dto.CreateGameDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.CreateGame(c.Context(), &dto); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "game created successfully"})
}

func (h *GameConfigHandler) UpdateSettings(c *fiber.Ctx) error {
	gameName := c.Params("game_name")

	patch, err := jsonpatch.DecodePatch(c.Body())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON Patch format"})
	}

	if err := h.service.UpdateSettings(c.Context(), gameName, patch); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "settings updated successfully"})
}

// =====================================================================
func (h *GameConfigHandler) UploadAsset(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	name, err := h.fileService.Upload(fileData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"url": "https://localhost:6069/api/game/assets/" + name})
}

// ==========================================================================
func (h *GameConfigHandler) GetAsset(c *fiber.Ctx) error {
	name := c.Params("name")
	file, err := h.fileService.Get(name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Disposition", "attachment; filename="+file.Name)
	c.Set("Content-Type", "application/octet-stream")

	contentType := h.fileService.GetContentType(file.Name)
	c.Set("Content-Type", contentType)

	return c.Send(file.Data)
}

func (h *GameConfigHandler) OverwriteSettings(c *fiber.Ctx) error {
	gameName := c.Params("game_name")
	var newSettings primitive.M

	if err := c.BodyParser(&newSettings); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid settings JSON"})
	}

	// Вызываем новый метод сервиса
	if err := h.service.OverwriteSettings(c.Context(), gameName, newSettings); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "settings updated successfully"})
}
