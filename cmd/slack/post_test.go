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

func TestPostCmd(t *testing.T) {
	// a mock slack client
	mockClient := &MockSlackClient{
		PostMessageFunc: func(channelID, text string) (string, error) {
			assert.Equal(t, "C1234567890", channelID)
			assert.Equal(t, "Hello, world!", text)
			return "1234567890.123456", nil
		},
	}

	// Override the NewSlackClient function
	oldNewSlackClient := NewSlackClient
	NewSlackClient = func(token string) slack.SlackClient {
		return mockClient
	}
	defer func() { NewSlackClient = oldNewSlackClient }()

	// Set the app token
	viper.Set("slack.app.token", "test-token")

	// Redirect STDIN
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.Write([]byte("Hello, world!"))
	w.Close()

	// Redirect STDOUT
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	r, w, _ = os.Pipe()
	os.Stdout = w

	// Run the command
	PostCmd.SetArgs([]string{"--channel", "C1234567890"})
	err := PostCmd.Execute()
	assert.NoError(t, err)

	// Check the output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Equal(t, "Message sent to Slack successfully with timestamp: 1234567890.123456\n", buf.String())
}