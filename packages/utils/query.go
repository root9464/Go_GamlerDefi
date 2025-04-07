package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Get[T any](url string) (T, error) {
	var result T

	status, body, errs := fiber.Get(url).Bytes()
	if len(errs) > 0 {
		return result, fmt.Errorf("request failed: %v", errs)
	}

	if status >= 400 {
		return result, fmt.Errorf("API error: %s", parseError(body, status))
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("decode error: %w", err)
	}

	return result, nil
}

func parseError(body []byte, status int) string {
	if len(body) == 0 {
		return fmt.Sprintf("status %d", status)
	}

	var data map[string]any
	if json.Unmarshal(body, &data) != nil {
		return fmt.Sprintf("%s (status %d)", body, status)
	}

	for _, key := range []string{"message", "error", "detail"} {
		if msg, ok := data[key].(string); ok {
			return fmt.Sprintf("%s (status %d)", msg, status)
		}
	}

	return fmt.Sprintf("%s (status %d)", body, status)
}
