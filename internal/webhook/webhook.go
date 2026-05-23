package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload represents the JSON body sent to the webhook endpoint.
type Payload struct {
	Job       string    `json:"job"`
	StartedAt time.Time `json:"started_at"`
	Duration  float64   `json:"duration_seconds"`
	ExitCode  int       `json:"exit_code"`
	Output    string    `json:"output"`
	Success   bool      `json:"success"`
}

// Client sends webhook notifications.
type Client struct {
	URL        string
	HTTPClient *http.Client
}

// New creates a new webhook Client with the given URL.
func New(url string) *Client {
	return &Client{
		URL: url,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send marshals the payload and POSTs it to the webhook URL.
func (c *Client) Send(p Payload) error {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := c.HTTPClient.Post(c.URL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("webhook: post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
