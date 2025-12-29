package ahasend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

type Client struct {
	httpClient *http.Client
	config     Config
	mu         sync.RWMutex
}

type Address struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type MessageRequest struct {
	From        Address   `json:"from"`
	Recipients  []Address `json:"recipients"`
	Subject     string    `json:"subject"`
	TextContent string    `json:"text_content,omitempty"`
	HTMLContent string    `json:"html_content,omitempty"`
}

func NewClient(config Config) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: config.Timeout},
		config:     config,
	}
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
}

func (c *Client) SendMessage(ctx context.Context, message MessageRequest) error {
	if len(message.Recipients) == 0 {
		return fmt.Errorf("ahasend: recipients required")
	}
	if strings.TrimSpace(message.TextContent) == "" && strings.TrimSpace(message.HTMLContent) == "" {
		return fmt.Errorf("ahasend: either text_content or html_content is required")
	}

	url := fmt.Sprintf("%s/accounts/%s/messages", strings.TrimRight(c.config.BaseURL, "/"), c.config.AccountID)
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("ahasend: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("ahasend: failed to create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ahasend: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		log.Error().
			Int("status", resp.StatusCode).
			Str("url", url).
			Str("response", string(respBody)).
			Msg("ahasend: request failed")
		return fmt.Errorf("ahasend: send failed with status %d", resp.StatusCode)
	}

	return nil
}
