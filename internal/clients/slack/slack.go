package slack

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

// Client is an interface that defines the methods for interacting with the Slack API.
type Client interface {
	PostMessage(channelID, text string) (string, error)
	DeleteMessage(channelID, timestamp string) error
	GetChannelID(channelName string) (string, error)
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

// GetChannelID retrieves the ID of a channel given its name.
func (c *client) GetChannelID(channelName string) (string, error) {
	// TODO: Implement pagination if there are more than 1000 channels.
	channels, _, err := c.api.GetConversations(&slack.GetConversationsParameters{
		Limit: 1000,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get conversations: %w", err)
	}

	// Normalize channel name for case-insensitive comparison.
	normalizedChannelName := strings.TrimPrefix(strings.ToLower(channelName), "#")

	for _, channel := range channels {
		if strings.ToLower(channel.Name) == normalizedChannelName {
			return channel.ID, nil
		}
	}

	return "", fmt.Errorf("channel '%s' not found", channelName)
}

