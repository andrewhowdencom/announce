package slack

// MockClient is a mock implementation of the Client interface for testing.
type MockClient struct {
	PostMessageFunc   func(channel, subject, text string) (string, error)
	DeleteMessageFunc func(channel, timestamp string) error
	GetChannelIDFunc  func(channelName string) (string, error)
}

// NewMockClient creates a new MockClient.
func NewMockClient() *MockClient {
	return &MockClient{
		PostMessageFunc: func(channel, subject, text string) (string, error) {
			return "1234567890.123456", nil
		},
		DeleteMessageFunc: func(channel, timestamp string) error {
			return nil
		},
		GetChannelIDFunc: func(channelName string) (string, error) {
			return "C1234567890", nil
		},
	}
}

// PostMessage calls the PostMessageFunc.
func (m *MockClient) PostMessage(channel, subject, text string) (string, error) {
	return m.PostMessageFunc(channel, subject, text)
}

// DeleteMessage calls the DeleteMessageFunc.
func (m *MockClient) DeleteMessage(channel, timestamp string) error {
	return m.DeleteMessageFunc(channel, timestamp)
}

// GetChannelID calls the GetChannelIDFunc.
func (m *MockClient) GetChannelID(channelName string) (string, error) {
	return m.GetChannelIDFunc(channelName)
}
