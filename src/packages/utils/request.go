package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type APIClient struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	Headers    http.Header
}

type RequestConfig struct {
	Ctx    context.Context
	Method string
	Path   string

	// Необязательные поля
	Headers     http.Header
	QueryParams url.Values
	Body        any
	Result      any
	ErrorResult any
}

func NewAPIClient(baseURL string) (*APIClient, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("некорректный базовый URL: %w", err)
	}

	return &APIClient{
		BaseURL: parsedBaseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 15,
		},
		Headers: make(http.Header),
	}, nil
}

func (c *APIClient) Do(config *RequestConfig) (*http.Response, error) {
	reqURL, err := c.BaseURL.Parse(config.Path)
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать URL запроса: %w", err)
	}
	if config.QueryParams != nil {
		reqURL.RawQuery = config.QueryParams.Encode()
	}

	var reqBody io.Reader
	if config.Body != nil {
		bodyBytes, err := json.Marshal(config.Body)
		if err != nil {
			return nil, fmt.Errorf("не удалось сериализовать тело запроса в JSON: %w", err)
		}
		reqBody = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequestWithContext(config.Ctx, config.Method, reqURL.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать HTTP-запрос: %w", err)
	}

	for key, values := range c.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	for key, values := range config.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	if config.Body != nil {
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
		}
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, fmt.Errorf("не удалось прочитать тело ответа: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		if config.Result != nil {
			if err := json.Unmarshal(respBody, config.Result); err != nil {
				return resp, fmt.Errorf("не удалось десериализовать успешный ответ: %w", err)
			}
		}
	} else {
		if config.ErrorResult != nil {
			if err := json.Unmarshal(respBody, config.ErrorResult); err != nil {
				return resp, fmt.Errorf("API вернул ошибку (статус %d): %s", resp.StatusCode, string(respBody))
			}
		}
		return resp, fmt.Errorf("API вернул ошибку (статус %d): %s", resp.StatusCode, string(respBody))
	}

	return resp, nil
}
