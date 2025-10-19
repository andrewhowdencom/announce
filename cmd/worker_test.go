package cmd

import (
	"testing"
	"time"

	"github.com/andrewhowdencom/announce/internal/clients/slack"
	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRunWorker(t *testing.T) {
	store, err := datastore.NewMockStore()
	assert.NoError(t, err)

	announcement := &datastore.Announcement{
		ID:          "1",
		Content:     "test",
		ChannelID:   "test",
		Status:      datastore.StatusPending,
		ScheduledAt: time.Now().Add(-1 * time.Hour),
	}
	err = store.AddAnnouncement(announcement)
	assert.NoError(t, err)

	slackClient := &slack.MockClient{
		PostMessageFunc: func(channelID, text string) (string, error) {
			return "12345", nil
		},
	}

	viper.Set("slack.app.token", "test")

	err = runWorker(store, slackClient)
	assert.NoError(t, err)

	updatedAnnouncement, err := store.GetAnnouncement("1")
	assert.NoError(t, err)

	assert.Equal(t, datastore.StatusProcessed, updatedAnnouncement.Status)

	sentMessages, err := store.ListSentMessagesByAnnouncementID("1")
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 1)
}
