// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package saavn

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseAPIURL  = "https://www.jiosaavn.com/api.php"
	DefaultUA   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	DefaultTime = 15 * time.Second
)

// Client handles communication with the JioSaavn API
type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
}

// NewClient creates a new JioSaavn HTTP client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: DefaultTime,
		},
		baseURL:   BaseAPIURL,
		userAgent: DefaultUA,
	}
}

// Request executes an API call with default Saavn query parameters
func (c *Client) Request(ctx context.Context, callName string, queryParams map[string]string) ([]byte, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	q := u.Query()
	q.Set("_format", "json")
	q.Set("_marker", "0")
	q.Set("api_version", "4")
	q.Set("ctx", "web6dot0")
	q.Set("__call", callName)

	for k, v := range queryParams {
		if v != "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("upstream returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
