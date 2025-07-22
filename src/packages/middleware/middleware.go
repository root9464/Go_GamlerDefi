package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	errors "github.com/root9464/Go_GamlerDefi/src/packages/lib/error"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

func LoggerMiddleware(logger *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Infof(
			"Request %s %s | IP: %s | User-Agent: %s",
			c.Method(),
			c.Path(),
			c.IP(),
			c.Context().UserAgent(),
		)

		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		logger.Infof(
			"Request completed | Duration: %.2fms | Status: %d",
			float64(latency.Microseconds())/1000,
			c.Response().StatusCode(),
		)

		if err != nil {
			logger.Errorf("Error: %v", err)
		}
		return err
	}

}

func ErrorMiddleware(ctx *fiber.Ctx) error {
	err := ctx.Next()

	if err != nil {
		if error, ok := err.(*errors.MapError); ok {
			return ctx.Status(error.Code).JSON(errors.MapError{
				Message:     error.Message,
				Description: error.Description,
			})
		}
		return ctx.Status(500).JSON(fiber.Map{
			"message": "internal server error",
		})
	}

	return nil
}

func SlicePrefix(prefix string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		originalPath := c.Path()
		if strings.HasPrefix(originalPath, prefix) {
			newPath := strings.TrimPrefix(originalPath, prefix)
			c.Path(newPath)
		}
		return c.Next()
	}
}
