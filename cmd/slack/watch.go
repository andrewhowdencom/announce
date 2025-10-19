/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package slack

import (
	"errors"
	"fmt"
	"time"

	"github.com/andrewhowdencom/announce/internal/clients/slack"
	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// ErrListAnnouncements is returned when the datastore fails to list announcements.
	ErrListAnnouncements = errors.New("failed to list announcements")
	// ErrUpdateAnnouncement is returned when the datastore fails to update an announcement.
	ErrUpdateAnnouncement = errors.New("failed to update announcement")
)

var (
	// newStore is a function that returns a new datastore.Storer.
	// It is defined as a variable so it can be replaced in tests.
	newStore = datastore.NewStore

	// newSlackClient is a function that returns a new slack.Client.
	// It is defined as a variable so it can be replaced in tests.
	newSlackClient = slack.NewClient
)

func doWatch(store datastore.Storer, slackClient slack.Client) error {
	announcements, err := store.ListAnnouncements()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrListAnnouncements, err)
	}

	for _, a := range announcements {
		if a.Status == datastore.StatusPending {
			fmt.Printf("Sending announcement %s to channel %s\n", a.ID, a.ChannelID)
			timestamp, err := slackClient.PostMessage(a.ChannelID, a.Content)
			if err != nil {
				a.Status = datastore.StatusFailed
			} else {
				a.Status = datastore.StatusSent
				a.Timestamp = timestamp
			}

			if err := store.UpdateAnnouncement(a); err != nil {
				return fmt.Errorf("%w: %w", ErrUpdateAnnouncement, err)
			}
		}
	}
	return nil
}

// WatchCmd represents the watch command
var WatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for pending announcements and send them",
	Long:  `Watch for pending announcements and send them.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := newStore()
		if err != nil {
			return fmt.Errorf("failed to create store: %w", err)
		}
		defer store.Close()

		slackToken := viper.GetString("slack.app.token")
		slackClient := newSlackClient(slackToken)

		fmt.Println("Watching for announcements...")
		for {
			if err := doWatch(store, slackClient); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			time.Sleep(1 * time.Second)
		}
	},
}
