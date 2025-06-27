package ton_controllers

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	ton_dto "github.com/root9464/Go_GamlerDefi/module/ton/dto"
)

// const (
// 	image_path = "../module/ton/files/logo.jpg"
// )

func (c *TonController) GetImage(ctx *fiber.Ctx) error {
	image_path := ctx.Params("image_path")
	c.logger.Infof("Image path: %s", image_path)

	baseDir, err := os.Getwd()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": "Oops, something went wrong"})
	}
	c.logger.Infof("Base dir: %s", baseDir)
	imagePath := filepath.Join(baseDir, "..", "module", "ton", "files", image_path)
	c.logger.Infof("Image path: %s", imagePath)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return ctx.Status(404).JSON(fiber.Map{"message": "Image not found"})
	}

	return ctx.SendFile(imagePath)
}

func (c *TonController) GetManifest(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(ton_dto.Manifest{
		URL:     "https://gamler.online",
		Name:    "Gamler",
		IconURL: "https://serv.gamler.online/web3/api/ton/image/logo.svg",
	})
}
