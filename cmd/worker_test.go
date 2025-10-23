package cmd

import (
	"testing"
	"time"

	"github.com/andrewhowdencom/ruf/internal/clients/slack"
	"github.com/andrewhowdencom/ruf/internal/datastore"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRunWorker(t *testing.T) {
	store, err := datastore.NewMockStore()
	assert.NoError(t, err)

	call := &datastore.Call{
		ID:          "1",
		Content:     "test",
		ChannelID:   "test",
		Status:      datastore.StatusPending,
		ScheduledAt: time.Now().Add(-1 * time.Hour),
	}
	err = store.AddCall(call)
	assert.NoError(t, err)

	slackClient := &slack.MockClient{
		PostMessageFunc: func(channelID, text string) (string, error) {
			return "12345", nil
		},
	}

	viper.Set("slack.app.token", "test")

	err = runWorker(store, slackClient)
	assert.NoError(t, err)

	updatedCall, err := store.GetCall("1")
	assert.NoError(t, err)

	assert.Equal(t, datastore.StatusProcessed, updatedCall.Status)

	sentMessages, err := store.ListSentMessagesByCallID("1")
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 1)
}
