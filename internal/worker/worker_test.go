package worker_test

import (
	"testing"
	"time"

	"github.com/andrewhowdencom/ruf/internal/clients/email"
	"github.com/andrewhowdencom/ruf/internal/clients/slack"
	"github.com/andrewhowdencom/ruf/internal/datastore"
	"github.com/andrewhowdencom/ruf/internal/model"
	"github.com/andrewhowdencom/ruf/internal/poller"
	"github.com/andrewhowdencom/ruf/internal/worker"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// mockSourcer implements the sourcer.Sourcer interface for testing.
type mockSourcer struct {
	callsBySource map[string][]*model.Call
	err           error
}

func (m *mockSourcer) Source(url string) ([]*model.Call, string, error) {
	if m.err != nil {
		return nil, "", m.err
	}
	return m.callsBySource[url], "state", nil
}

func TestWorker_RunTick(t *testing.T) {
	// Mock datastore
	store := datastore.NewMockStore()

	// Mock Slack client
	slackClient := slack.NewMockClient()

	// Mock Email client
	emailClient := email.NewMockClient()

	// Mock sourcer
	s := &mockSourcer{
		callsBySource: map[string][]*model.Call{
			"mock://url": {
				{
					ID:      "1",
					Subject: "Test Subject",
					Content: "Hello, world!",
					Destinations: []model.Destination{
						{
							Type: "slack",
							To:   []string{"test-channel"},
						},
						{
							Type: "email",
							To:   []string{"test@example.com"},
						},
					},
					ScheduledAt: time.Now().Add(-1 * time.Minute),
				},
			},
		},
	}

	p := poller.New(s, 1*time.Minute)
	viper.Set("source.urls", []string{"mock://url"})

	w := worker.New(store, slackClient, emailClient, p, 1*time.Minute)

	err := w.RunTick()
	assert.NoError(t, err)

	sentMessages, err := store.ListSentMessages()
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 2)
	assert.Equal(t, "1", sentMessages[0].SourceID)
	assert.Equal(t, "1", sentMessages[1].SourceID)
}

func TestWorker_RunTickWithOldCall(t *testing.T) {
	// Mock datastore
	store := datastore.NewMockStore()

	// Mock Slack client
	slackClient := slack.NewMockClient()

	// Mock Email client
	emailClient := email.NewMockClient()

	// Mock sourcer
	s := &mockSourcer{
		callsBySource: map[string][]*model.Call{
			"mock://url": {
				{
					ID:      "1",
					Content: "Hello, world!",
					Destinations: []model.Destination{
						{
							Type: "slack",
							To:   []string{"test-channel"},
						},
					},
					ScheduledAt: time.Now().Add(-48 * time.Hour),
				},
			},
		},
	}

	p := poller.New(s, 1*time.Minute)

	viper.Set("source.urls", []string{"mock://url"})
	viper.Set("worker.lookback_period", "24h")

	w := worker.New(store, slackClient, emailClient, p, 1*time.Minute)

	err := w.RunTick()
	assert.NoError(t, err)

	sentMessages, err := store.ListSentMessages()
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 1)
	assert.Equal(t, "1", sentMessages[0].SourceID)
	assert.Equal(t, datastore.StatusFailed, sentMessages[0].Status)
}

func TestWorker_RunTickWithDeletedCall(t *testing.T) {
	// Mock datastore
	store := datastore.NewMockStore()

	// Mock Slack client
	slackClient := slack.NewMockClient()

	// Mock Email client
	emailClient := email.NewMockClient()

	scheduledAt := time.Now().Add(-1 * time.Minute).UTC()

	// Add a deleted message to the store
	err := store.AddSentMessage(&datastore.SentMessage{
		SourceID:    "1",
		ScheduledAt: scheduledAt,
		Status:      datastore.StatusDeleted,
		Type:        "slack",
		Destination: "test-channel",
	})
	assert.NoError(t, err)

	// Mock sourcer
	s := &mockSourcer{
		callsBySource: map[string][]*model.Call{
			"mock://url": {
				{
					ID:      "1",
					Subject: "Test Subject",
					Content: "Hello, world!",
					Destinations: []model.Destination{
						{
							Type: "slack",
							To:   []string{"test-channel"},
						},
					},
					ScheduledAt: scheduledAt,
				},
			},
		},
	}

	p := poller.New(s, 1*time.Minute)

	viper.Set("source.urls", []string{"mock://url"})

	w := worker.New(store, slackClient, emailClient, p, 1*time.Minute)

	err = w.RunTick()
	assert.NoError(t, err)

	// Check that the slack client was not called
	assert.Equal(t, 0, slackClient.PostMessageCount)
}
