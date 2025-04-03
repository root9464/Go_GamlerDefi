package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
)

func LoggerMiddleware(logger *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)
		latencyStr := fmt.Sprintf("%.2fms", float64(latency.Microseconds())/1000)
		logger.Infof(
			"Request %s %s | IP: %s | User-Agent: %s | Duration: %s | Status: %d",
			c.Method(),
			c.Path(),
			c.IP(),
			c.Context().UserAgent(),
			latencyStr,
			c.Response().StatusCode(),
		)
		if err != nil {
			logger.Errorf("Error: %v", err)
		}
		return err
	}
}
