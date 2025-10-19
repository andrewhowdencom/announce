package slack

import (
	"fmt"
	"testing"

	"github.com/andrewhowdencom/announce/internal/clients/slack"
	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDoDelete_Sent(t *testing.T) {
	a := &datastore.Announcement{
		ID:        "1",
		Content:   "test",
		ChannelID: "test",
		Status:    datastore.StatusSent,
	}

	deleteAnnouncementCalled := false
	deleteSentMessageCalled := false
	deleteMessageCalled := false

	store := &datastore.MockStore{
		GetAnnouncementFunc: func(id string) (*datastore.Announcement, error) {
			if id == "1" {
				return a, nil
			}
			return nil, fmt.Errorf("announcement not found")
		},
		DeleteAnnouncementFunc: func(id string) error {
			deleteAnnouncementCalled = true
			return nil
		},
		ListSentMessagesByAnnouncementIDFunc: func(announcementID string) ([]*datastore.SentMessage, error) {
			return []*datastore.SentMessage{
				{
					ID:             "1",
					AnnouncementID: "1",
					Timestamp:      "12345",
					Status:         datastore.StatusSent,
				},
			}, nil
		},
		DeleteSentMessageFunc: func(id string) error {
			deleteSentMessageCalled = true
			return nil
		},
	}

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

	err := doDelete(store, client, "1")
	assert.NoError(t, err)

	assert.True(t, deleteAnnouncementCalled)
	assert.True(t, deleteSentMessageCalled)
	assert.True(t, deleteMessageCalled)
}

func TestDoDelete_Pending(t *testing.T) {
	a := &datastore.Announcement{
		ID:        "1",
		Content:   "test",
		ChannelID: "test",
		Status:    datastore.StatusPending,
	}

	deleteAnnouncementCalled := false
	deleteSentMessageCalled := false
	deleteMessageCalled := false

	store := &datastore.MockStore{
		GetAnnouncementFunc: func(id string) (*datastore.Announcement, error) {
			if id == "1" {
				return a, nil
			}
			return nil, fmt.Errorf("announcement not found")
		},
		DeleteAnnouncementFunc: func(id string) error {
			deleteAnnouncementCalled = true
			return nil
		},
		ListSentMessagesByAnnouncementIDFunc: func(announcementID string) ([]*datastore.SentMessage, error) {
			return []*datastore.SentMessage{}, nil
		},
		DeleteSentMessageFunc: func(id string) error {
			deleteSentMessageCalled = true
			return nil
		},
	}

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

	err := doDelete(store, client, "1")
	assert.NoError(t, err)

	assert.True(t, deleteAnnouncementCalled)
	assert.False(t, deleteSentMessageCalled)
	assert.False(t, deleteMessageCalled)
}
