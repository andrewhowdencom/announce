/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package slack

import "github.com/andrewhowdencom/announce/internal/clients/slack"

// NewSlackClient is a function that returns a new Slack client. It can be
// replaced in tests with a mock client.
var NewSlackClient = func(token string) slack.SlackClient {
	return slack.NewClient(token)
}

// MockSlackClient is a mock implementation of the SlackClient interface
type MockSlackClient struct {
	PostMessageFunc   func(channelID, text string) (string, error)
	DeleteMessageFunc func(channelID, timestamp string) error
}

func (m *MockSlackClient) PostMessage(channelID, text string) (string, error) {
	return m.PostMessageFunc(channelID, text)
}

func (m *MockSlackClient) DeleteMessage(channelID, timestamp string) error {
	return m.DeleteMessageFunc(channelID, timestamp)
}