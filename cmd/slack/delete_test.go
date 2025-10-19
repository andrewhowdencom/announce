package slack

import (
	"testing"

	"github.com/andrewhowdencom/announce/internal/clients/slack"
	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDoDelete_Sent(t *testing.T) {
	store, err := datastore.NewMockStore()
	assert.NoError(t, err)

	a := &datastore.Announcement{
		ID:        "1",
		Content:   "test",
		ChannelID: "test",
		Status:    datastore.StatusSent,
	}
	err = store.AddAnnouncement(a)
	assert.NoError(t, err)

	sm := &datastore.SentMessage{
		ID:             "1",
		AnnouncementID: "1",
		Timestamp:      "12345",
		Status:         datastore.StatusSent,
	}
	err = store.AddSentMessage(sm)
	assert.NoError(t, err)

	deleteMessageCalled := false
	client := &slack.MockClient{
		GetChannelIDFunc: func(channelName string) (string, error) {
			return "test", nil
		},
		DeleteMessageFunc: func(channelID, timestamp string) error {
			deleteMessageCalled = true
			return nil
		},
	}
	viper.Set("slack.app.token", "test")

	err = doDelete(store, client, "1")
	assert.NoError(t, err)

	_, err = store.GetAnnouncement("1")
	assert.ErrorIs(t, err, datastore.ErrNotFound)

	sentMessages, err := store.ListSentMessagesByAnnouncementID("1")
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 0)

	assert.True(t, deleteMessageCalled)
}

func TestDoDelete_Pending(t *testing.T) {
	store, err := datastore.NewMockStore()
	assert.NoError(t, err)

	a := &datastore.Announcement{
		ID:        "1",
		Content:   "test",
		ChannelID: "test",
		Status:    datastore.StatusPending,
	}
	err = store.AddAnnouncement(a)
	assert.NoError(t, err)

	deleteMessageCalled := false
	client := &slack.MockClient{
		GetChannelIDFunc: func(channelName string) (string, error) {
			return "test", nil
		},
		DeleteMessageFunc: func(channelID, timestamp string) error {
			deleteMessageCalled = true
			return nil
		},
	}
	viper.Set("slack.app.token", "test")

	err = doDelete(store, client, "1")
	assert.NoError(t, err)

	_, err = store.GetAnnouncement("1")
	assert.ErrorIs(t, err, datastore.ErrNotFound)

	assert.False(t, deleteMessageCalled)
}
