package slack

import (
	"github.com/slack-go/slack"
)

// SlackClient defines the interface for interacting with Slack.
type SlackClient interface {
	PostMessage(channelID, text string) (string, error)
	DeleteMessage(channelID, timestamp string) error
}

// client is a wrapper around the slack client that implements SlackClient
type client struct {
	api *slack.Client
}

// NewClient creates a new slack client
func NewClient(token string) SlackClient {
	api := slack.New(token)
	return &client{api: api}
}

// PostMessage sends a message to a slack channel
func (c *client) PostMessage(channelID, text string) (string, error) {
	// The Slack API returns the channel ID, timestamp, and an error.
	// We are only interested in the timestamp, which is needed for deleting messages.
	_, timestamp, err := c.api.PostMessage(channelID, slack.MsgOptionText(text, false))
	if err != nil {
		return "", err
	}
	return timestamp, nil
}

// DeleteMessage deletes a message from a slack channel
func (c *client) DeleteMessage(channelID, timestamp string) error {
	_, _, err := c.api.DeleteMessage(channelID, timestamp)
	return err
}