package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client represents a client for interacting with the Mailer service.
type Client struct {
	server string
	client *http.Client
}

// NewClient creates a new Client instance.
func NewClient(serverAddr string, client ...*http.Client) *Client {
	var cli *http.Client
	if len(client) > 0 {
		cli = client[0]
	}

	if cli == nil {
		cli = http.DefaultClient
	}

	return &Client{
		server: serverAddr,
		client: cli,
	}
}

// SendText sends a text message.
func (c *Client) SendText(to, subject, text string) error {
	return c.send(&Message{
		ToAddress: to,
		Subject:   subject,
		Content: &Content{
			Body: text,
			HTML: false,
		},
	})
}

// SendHTML sends an HTML message.
func (c *Client) SendHTML(to, subject, html string) error {
	return c.send(&Message{
		ToAddress: to,
		Subject:   subject,
		Content: &Content{
			Body: html,
			HTML: true,
		},
	})
}

// SendTemplate sends a template message.
func (c *Client) SendTemplate(to, subject string, template TemplateID, props TemplateProps) error {
	return c.send(&Message{
		ToAddress: to,
		Subject:   subject,
		Template: &Template{
			ID:    template.String(),
			Props: props,
		},
	})
}

// send sends the Message
func (c *Client) send(msg *Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("Failed to marshal message: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.server+"/send", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("Failed to create request: %v", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to send request: %v", err)
	}
	defer res.Body.Close()

	b, _ = io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		unexpErr := fmt.Errorf("unexpected status code: %d", res.StatusCode)
		if b != nil {
			return fmt.Errorf("%s: %w", string(b), unexpErr)
		}
		return fmt.Errorf("Unknown error: %w", unexpErr)
	}

	return nil
}
