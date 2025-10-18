/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package slack

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/andrewhowdencom/announce/internal/clients/slack"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCmd(t *testing.T) {
	// a mock slack client
	mockClient := &MockSlackClient{
		DeleteMessageFunc: func(channelID, timestamp string) error {
			assert.Equal(t, "C1234567890", channelID)
			assert.Equal(t, "1234567890.123456", timestamp)
			return nil
		},
	}

	// Override the NewSlackClient function
	oldNewSlackClient := NewSlackClient
	NewSlackClient = func(token string) slack.SlackClient {
		return mockClient
	}
	defer func() { NewSlackClient = oldNewSlackClient }()

	// Set the bot token
	viper.Set("slack.bot_token", "test-token")

	// Redirect STDOUT
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	DeleteCmd.SetArgs([]string{"--channel", "C1234567890", "--timestamp", "1234567890.123456"})
	err := DeleteCmd.Execute()
	assert.NoError(t, err)

	// Check the output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Equal(t, "Message deleted from Slack successfully\n", buf.String())
}