/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package slack

import (
	"testing"

	"github.com/andrewhowdencom/announce/internal/clients/slack"
	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/spf13/viper"
)

func TestDoWatch(t *testing.T) {
	announcements := []*datastore.Announcement{
		{
			ID:        "1",
			Content:   "test",
			ChannelID: "test",
			Status:    datastore.StatusPending,
		},
	}

	store := &datastore.MockStore{
		ListAnnouncementsFunc: func() ([]*datastore.Announcement, error) {
			return announcements, nil
		},
		UpdateAnnouncementFunc: func(a *datastore.Announcement) error {
			announcements[0] = a
			return nil
		},
	}

	slackClient := &slack.MockClient{
		PostMessageFunc: func(channelID, text string) (string, error) {
			return "12345", nil
		},
	}

	viper.Set("slack.app.token", "test")

	if err := doWatch(store, slackClient); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if announcements[0].Status != datastore.StatusSent {
		t.Errorf("expected status to be sent, but got %s", announcements[0].Status)
	}
}
