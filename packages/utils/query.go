package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

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

func Patch[T any](url string, body any) (T, error) {
	var result T

	agent := fiber.Patch(url)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return result, fmt.Errorf("encode request body error: %w", err)
	}

	agent.ContentType("application/json")
	agent.Body(jsonBody)

	status, respBody, errs := agent.Bytes()
	if len(errs) > 0 {
		return result, fmt.Errorf("request failed: %v", errs)
	}

	if status >= 400 {
		return result, fmt.Errorf("API error: %s", parseError(respBody, status))
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return result, fmt.Errorf("decode response error: %w", err)
	}

	return result, nil
}
