package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPRequestInput represents the input schema for the HTTP request tool
type HTTPRequestInput struct {
	URL     string            `json:"url"`
	Method  string           `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string           `json:"body,omitempty"`
}

// HTTPRequestOutput represents the output schema for the HTTP request tool
type HTTPRequestOutput struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string           `json:"body"`
}

// NewHTTPRequestTool creates a new HTTP request tool
func NewHTTPRequestTool() (string, string, json.RawMessage, func(context.Context, json.RawMessage) (json.RawMessage, error)) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"url": map[string]interface{}{
				"type": "string",
				"description": "The URL to send the request to",
			},
			"method": map[string]interface{}{
				"type": "string",
				"enum": []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"},
				"description": "The HTTP method to use",
			},
			"headers": map[string]interface{}{
				"type": "object",
				"additionalProperties": map[string]interface{}{
					"type": "string",
				},
				"description": "Optional headers to include in the request",
			},
			"body": map[string]interface{}{
				"type": "string",
				"description": "Optional body to include in the request",
			},
		},
		"required": []string{"url", "method"},
	}

	schemaJSON, _ := json.Marshal(schema)

	return "httpRequest",
		"Makes HTTP requests to external services",
		schemaJSON,
		func(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
			var params HTTPRequestInput
			if err := json.Unmarshal(input, &params); err != nil {
				return nil, fmt.Errorf("invalid input: %w", err)
			}

			// Create HTTP client with timeout
			client := &http.Client{
				Timeout: 30 * time.Second,
			}

			// Create request
			req, err := http.NewRequestWithContext(ctx, params.Method, params.URL, strings.NewReader(params.Body))
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}

			// Add headers
			for key, value := range params.Headers {
				req.Header.Add(key, value)
			}

			// Execute request
			resp, err := client.Do(req)
			if err != nil {
				return nil, fmt.Errorf("request failed: %w", err)
			}
			defer resp.Body.Close()

			// Read response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %w", err)
			}

			// Convert response headers
			headers := make(map[string]string)
			for key, values := range resp.Header {
				headers[key] = strings.Join(values, ", ")
			}

			// Create output
			output := HTTPRequestOutput{
				StatusCode: resp.StatusCode,
				Headers:    headers,
				Body:      string(body),
			}

			outputJSON, err := json.Marshal(output)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal output: %w", err)
			}

			return outputJSON, nil
		}
} 