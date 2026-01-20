package vast

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://console.vast.ai/api/v0"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchOffers executes the search query
func (c *Client) SearchOffers(ctx context.Context, builder *SearchBuilder) ([]Offer, error) {
	endpoint := baseURL + "/bundles/"

	// 1. Build JSON body
	payload, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	// 2. Create Request
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// 3. Execute
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// エラーハンドリング（本番ではBodyを読んで詳細なエラーを返すのが望ましい）
		return nil, fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	// 4. Decode
	var result searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Offers, nil
}
