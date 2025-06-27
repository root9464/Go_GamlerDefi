package ton_controllers

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

const (
	image_path = "../files/logo.jpg"
)

func (c *TonController) GetImage(ctx *fiber.Ctx) error {
	if _, err := os.Stat(image_path); os.IsNotExist(err) {
		return ctx.Status(404).SendString("Image not found")
	}
	return ctx.SendFile(image_path)
}
