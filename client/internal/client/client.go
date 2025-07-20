package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Register(ctx context.Context, username, password string) error {
	req := map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := c.client.Post(
		c.baseURL+"/register",
		"application/json",
		bytes.NewBuffer(mustJson(req)),
	)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return parseError(resp)
	}
	return nil
}

func (c *Client) Login(ctx context.Context, username, password string) (string, error) {
	req := map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := c.client.Post(
		c.baseURL+"/login",
		"application/json",
		bytes.NewBuffer(mustJson(req)),
	)
	if err != nil {
		return "", fmt.Errorf("connection error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", parseError(resp)
	}

	var result struct{ Token string }
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Token, nil
}

func mustJson(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func parseError(resp *http.Response) error {
	var errResp struct{ Error string }
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}
	return fmt.Errorf(errResp.Error)
}
