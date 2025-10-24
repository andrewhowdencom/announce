package cmd

import (
	"testing"
	"time"

	"github.com/andrewhowdencom/ruf/internal/clients/slack"
	"github.com/andrewhowdencom/ruf/internal/datastore"
	"github.com/andrewhowdencom/ruf/internal/model"
	"github.com/andrewhowdencom/ruf/internal/worker"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type mockSourcer struct {
	calls []*model.Call
	err   error
}

func (m *mockSourcer) Source(url string) ([]*model.Call, error) {
	return m.calls, m.err
}

type mockEmailClient struct {
	sendFunc func(to []string, subject, body string) error
}

func (m *mockEmailClient) Send(to []string, subject, body string) error {
	if m.sendFunc != nil {
		return m.sendFunc(to, subject, body)
	}
	return nil
}

func TestRunTick(t *testing.T) {
	t.Run("slack call", func(t *testing.T) {
		// Mock datastore
		store := datastore.NewMockStore()

		// Mock Slack client
		slackClient := slack.NewMockClient()

		// Mock email worker
		emailWorker := worker.NewEmailWorker(&mockEmailClient{})

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

		viper.Set("source.urls", []string{"mock://url"})

		err := runTick(store, slackClient, emailWorker, s)
		assert.NoError(t, err)

		sentMessages, err := store.ListSentMessages()
		assert.NoError(t, err)
		assert.Len(t, sentMessages, 1)
		assert.Equal(t, "1", sentMessages[0].SourceID)
	})

	t.Run("email call", func(t *testing.T) {
		// Mock datastore
		store := datastore.NewMockStore()

		// Mock Slack client
		slackClient := slack.NewMockClient()

		// Mock email worker
		var sent bool
		emailWorker := worker.NewEmailWorker(&mockEmailClient{
			sendFunc: func(to []string, subject, body string) error {
				sent = true
				return nil
			},
		})

		// Mock sourcer
		s := &mockSourcer{
			calls: []*model.Call{
				{
					ID: "2",
					Email: &model.Email{
						To:      []string{"test@example.com"},
						Subject: "Test",
						Body:    "Test body",
					},
					ScheduledAt: time.Now().Add(-1 * time.Minute),
				},
			},
		}

		viper.Set("source.urls", []string{"mock://url"})

		err := runTick(store, slackClient, emailWorker, s)
		assert.NoError(t, err)

		sentMessages, err := store.ListSentMessages()
		assert.NoError(t, err)
		assert.Len(t, sentMessages, 1)
		assert.Equal(t, "2", sentMessages[0].SourceID)
		assert.True(t, sent)
	})
}
