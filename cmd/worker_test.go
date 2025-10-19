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
	announcements := []*datastore.Announcement{
		{
			ID:          "1",
			Content:     "test",
			ChannelID:   "test",
			Status:      datastore.StatusPending,
			ScheduledAt: time.Now().Add(-1 * time.Hour),
		},
	}

	addSentMessageCalled := false

	store := &datastore.MockStore{
		ListAnnouncementsFunc: func() ([]*datastore.Announcement, error) {
			return announcements, nil
		},
		UpdateAnnouncementFunc: func(a *datastore.Announcement) error {
			announcements[0] = a
			return nil
		},
		AddSentMessageFunc: func(sm *datastore.SentMessage) error {
			addSentMessageCalled = true
			return nil
		},
	}

	slackClient := &slack.MockClient{
		PostMessageFunc: func(channelID, text string) (string, error) {
			return "12345", nil
		},
	}

	viper.Set("slack.app.token", "test")

	err := runWorker(store, slackClient)
	assert.NoError(t, err)

	assert.Equal(t, datastore.StatusProcessed, announcements[0].Status)
	assert.True(t, addSentMessageCalled)
}
