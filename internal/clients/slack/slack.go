package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

// Client is an interface that defines the methods for interacting with the Slack API.
type Client interface {
	PostMessage(channelID, text string) (string, error)
	DeleteMessage(channelID, timestamp string) error
}

// client is the concrete implementation of the Client interface.
type client struct {
	api *slack.Client
}

// NewClient creates a new Slack client.
func NewClient(token string) Client {
	return &client{
		api: slack.New(token),
	}
}

// PostMessage sends a message to a Slack channel.
func (c *client) PostMessage(channelID, text string) (string, error) {
	_, timestamp, err := c.api.PostMessage(channelID, slack.MsgOptionText(text, false))
	if err != nil {
		return "", fmt.Errorf("failed to post message: %w", err)
	}
	return timestamp, nil
}

// DeleteMessage deletes a message from a Slack channel.
func (c *client) DeleteMessage(channelID, timestamp string) error {
	_, _, err := c.api.DeleteMessage(channelID, timestamp)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}

// MockClient is a mock implementation of the Client interface for testing.
type MockClient struct {
	PostMessageFunc   func(channelID, text string) (string, error)
	DeleteMessageFunc func(channelID, timestamp string) error
}

// PostMessage calls the PostMessageFunc.
func (m *MockClient) PostMessage(channelID, text string) (string, error) {
	return m.PostMessageFunc(channelID, text)
}

// DeleteMessage calls the DeleteMessageFunc.
func (m *MockClient) DeleteMessage(channelID, timestamp string) error {
	return m.DeleteMessageFunc(channelID, timestamp)
}
