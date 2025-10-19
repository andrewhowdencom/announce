package slack

import (
	"fmt"
	"testing"

	"github.com/andrewhowdencom/announce/internal/clients/slack"
	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/spf13/viper"
)

func TestDoDelete_SoftDelete(t *testing.T) {
	a := &datastore.Announcement{
		ID:        "1",
		Content:   "test",
		ChannelID: "test",
		Status:    datastore.StatusSent,
		Timestamp: "12345",
	}

	store := &datastore.MockStore{
		GetAnnouncementFunc: func(id string) (*datastore.Announcement, error) {
			if id == "1" {
				return a, nil
			}
			return nil, fmt.Errorf("announcement not found")
		},
		UpdateAnnouncementFunc: func(updatedAnnouncement *datastore.Announcement) error {
			a = updatedAnnouncement
			return nil
		},
	}

	client := &slack.MockClient{
		GetChannelIDFunc: func(channelName string) (string, error) {
			return "test", nil
		},
		DeleteMessageFunc: func(channelID, timestamp string) error {
			return nil
		},
	}
	viper.Set("slack.app.token", "test")

	if err := doDelete(store, client, "1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a.Status != datastore.StatusDeleted {
		t.Errorf("expected status to be deleted, but got %s", a.Status)
	}
}

func TestDoDelete_HardDelete(t *testing.T) {
	a := &datastore.Announcement{
		ID:        "1",
		Content:   "test",
		ChannelID: "test",
		Status:    datastore.StatusPending,
	}
	deleted := false

	store := &datastore.MockStore{
		GetAnnouncementFunc: func(id string) (*datastore.Announcement, error) {
			if id == "1" {
				return a, nil
			}
			return nil, fmt.Errorf("announcement not found")
		},
		DeleteAnnouncementFunc: func(id string) error {
			if id == "1" {
				deleted = true
				return nil
			}
			return fmt.Errorf("announcement not found")
		},
	}
	viper.Set("slack.app.token", "test")

	client := &slack.MockClient{}

	if err := doDelete(store, client, "1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !deleted {
		t.Error("expected announcement to be deleted, but it wasn't")
	}
}

func TestDoDelete_ResolvesChannelName(t *testing.T) {
	a := &datastore.Announcement{
		ID:        "1",
		Content:   "test",
		ChannelID: "#general",
		Status:    datastore.StatusSent,
		Timestamp: "12345",
	}

	store := &datastore.MockStore{
		GetAnnouncementFunc: func(id string) (*datastore.Announcement, error) {
			if id == "1" {
				return a, nil
			}
			return nil, fmt.Errorf("announcement not found")
		},
		UpdateAnnouncementFunc: func(updatedAnnouncement *datastore.Announcement) error {
			a = updatedAnnouncement
			return nil
		},
	}

	client := &slack.MockClient{
		GetChannelIDFunc: func(channelName string) (string, error) {
			if channelName == "#general" {
				return "C1234567890", nil
			}
			return "", fmt.Errorf("channel not found")
		},
		DeleteMessageFunc: func(channelID, timestamp string) error {
			if channelID != "C1234567890" {
				return fmt.Errorf("expected channel ID to be C1234567890, but got %s", channelID)
			}
			return nil
		},
	}
	viper.Set("slack.app.token", "test")

	if err := doDelete(store, client, "1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a.Status != datastore.StatusDeleted {
		t.Errorf("expected status to be deleted, but got %s", a.Status)
	}
}
