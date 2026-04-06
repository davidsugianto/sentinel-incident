package lark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	webhookURL string
	httpClient *http.Client
}

// CardMessage represents Lark's interactive card message
type CardMessage struct {
	MsgType string `json:"msg_type"`
	Card    *Card  `json:"card,omitempty"`
}

// Card represents a Lark card
type Card struct {
	Config   *CardConfig `json:"config,omitempty"`
	Header   *CardHeader `json:"header,omitempty"`
	Elements []Element   `json:"elements,omitempty"`
}

// CardConfig for card behavior
type CardConfig struct {
	WideScreenMode bool `json:"wide_screen_mode"`
	EnableForward  bool `json:"enable_forward"`
}

// CardHeader for card title
type CardHeader struct {
	Title    *CardTitle `json:"title,omitempty"`
	Template string     `json:"template,omitempty"` // red, green, blue, etc.
}

// CardTitle for header title
type CardTitle struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

// Element interface for card elements (using interface{} for flexibility)
type Element interface{}

// DivElement for text content
type DivElement struct {
	Tag  string    `json:"tag"`
	Text *TextItem `json:"text,omitempty"`
}

// FieldElement for key-value pairs
type FieldElement struct {
	Tag    string  `json:"tag"`
	Fields []Field `json:"fields,omitempty"`
}

// Field for key-value display
type Field struct {
	IsShort bool     `json:"is_short"`
	Text    TextItem `json:"text"`
}

// TextItem for text content
type TextItem struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

func NewClient(webhookURL string) *Client {
	return &Client{
		webhookURL: webhookURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) SendMessage(ctx context.Context, msg *CardMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("lark webhook returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
