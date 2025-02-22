package wikipedia

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// WikipediaInput represents the input schema for the Wikipedia tool
type WikipediaInput struct {
	Query string `json:"query"`
}

// WikipediaOutput represents the output schema for the Wikipedia tool
type WikipediaOutput struct {
	Title    string `json:"title"`
	Extract  string `json:"extract"`
	URL      string `json:"url"`
	PageID   int    `json:"pageId"`
}

// wikipediaAPIResponse represents the Wikipedia API response
type wikipediaAPIResponse struct {
	Query struct {
		Pages map[string]struct {
			PageID  int    `json:"pageid"`
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

// NewWikipediaTool creates a new Wikipedia search tool
func NewWikipediaTool() (string, string, json.RawMessage, func(context.Context, json.RawMessage) (json.RawMessage, error)) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type": "string",
				"description": "The search query for Wikipedia",
			},
		},
		"required": []string{"query"},
	}

	schemaJSON, _ := json.Marshal(schema)

	return "wikipedia",
		"Searches Wikipedia for information about a topic",
		schemaJSON,
		func(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
			var params WikipediaInput
			if err := json.Unmarshal(input, &params); err != nil {
				return nil, fmt.Errorf("invalid input: %w", err)
			}

			// Create Wikipedia API URL
			apiURL := fmt.Sprintf(
				"https://en.wikipedia.org/w/api.php?action=query&format=json&prop=extracts&exintro=true&explaintext=true&titles=%s",
				url.QueryEscape(params.Query),
			)

			// Create request
			req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}

			// Add headers
			req.Header.Add("User-Agent", "Go-Tools-Agent/1.0")

			// Execute request
			client := &http.Client{}
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

			// Parse response
			var apiResp wikipediaAPIResponse
			if err := json.Unmarshal(body, &apiResp); err != nil {
				return nil, fmt.Errorf("failed to parse response: %w", err)
			}

			// Get first page from response
			var output WikipediaOutput
			for _, page := range apiResp.Query.Pages {
				output = WikipediaOutput{
					Title:    page.Title,
					Extract:  strings.TrimSpace(page.Extract),
					URL:      fmt.Sprintf("https://en.wikipedia.org/?curid=%d", page.PageID),
					PageID:   page.PageID,
				}
				break
			}

			if output.Title == "" {
				return nil, fmt.Errorf("no results found for query: %s", params.Query)
			}

			outputJSON, err := json.Marshal(output)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal output: %w", err)
			}

			return outputJSON, nil
		}
} 