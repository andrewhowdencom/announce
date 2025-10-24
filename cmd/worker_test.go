package cmd

import (
	"testing"
	"time"

	"github.com/andrewhowdencom/ruf/internal/clients/slack"
	"github.com/andrewhowdencom/ruf/internal/datastore"
	"github.com/andrewhowdencom/ruf/internal/model"
	"github.com/andrewhowdencom/ruf/internal/poller"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type mockSourcer struct {
	calls []*model.Call
	err   error
}

func (m *mockSourcer) Source(url string) ([]*model.Call, string, error) {
	return m.calls, "state", m.err
}

func TestRunTick(t *testing.T) {
	// Mock datastore
	store := datastore.NewMockStore()

	// Mock Slack client
	slackClient := slack.NewMockClient()

	// Mock sourcer
	s := &mockSourcer{
		calls: []*model.Call{
			{
				ID:      "1",
				Content: "Hello, world!",
				Destinations: []model.Destination{
					{
						Type:      "slack",
						ChannelID: "C1234567890",
					},
				},
				ScheduledAt: time.Now().Add(-1 * time.Minute),
			},
		},
	}

	p := poller.New(s, 1*time.Minute)

	viper.Set("source.urls", []string{"mock://url"})

	err := runTick(store, slackClient, p)
	assert.NoError(t, err)

	sentMessages, err := store.ListSentMessages()
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 1)
	assert.Equal(t, "1", sentMessages[0].SourceID)
}

func TestRunTickWithOldCall(t *testing.T) {
	// Mock datastore
	store := datastore.NewMockStore()

	// Mock Slack client
	slackClient := slack.NewMockClient()

	// Mock sourcer
	s := &mockSourcer{
		calls: []*model.Call{
			{
				ID:      "1",
				Content: "Hello, world!",
				Destinations: []model.Destination{
					{
						Type:      "slack",
						ChannelID: "C1234567890",
					},
				},
				ScheduledAt: time.Now().Add(-48 * time.Hour),
			},
		},
	}

	p := poller.New(s, 1*time.Minute)

	viper.Set("source.urls", []string{"mock://url"})
	viper.Set("worker.lookback_period", "24h")

	err := runTick(store, slackClient, p)
	assert.NoError(t, err)

	sentMessages, err := store.ListSentMessages()
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 1)
	assert.Equal(t, "1", sentMessages[0].SourceID)
	assert.Equal(t, datastore.StatusFailed, sentMessages[0].Status)
}
