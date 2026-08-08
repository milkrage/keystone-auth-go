package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func (c *Client) IssueToken(ctx context.Context, r Request) (string, time.Time, error) {
	path, err := url.JoinPath(c.BaseURL, "/v3/auth/tokens")
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to build url: %w", err)
	}

	body, err := json.Marshal(r)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "keystone-auth-go")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", time.Time{}, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	token := resp.Header.Get("X-Subject-Token")
	if token == "" {
		return "", time.Time{}, fmt.Errorf("response does not contain an x-subject-token header")
	}

	result := Response{}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return token, result.Token.ExpiresAt, nil
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}

	return http.DefaultClient
}
