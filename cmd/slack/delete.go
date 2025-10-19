package slack

import (
	"fmt"

	"github.com/andrewhowdencom/announce/internal/clients/slack"
	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete an announcement",
	Long:  `Delete an announcement.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := datastore.NewStore()
		if err != nil {
			return fmt.Errorf("failed to create datastore: %w", err)
		}
		defer store.Close()

		client := slack.NewClient(viper.GetString("slack.app.token"))

		return doDelete(store, client, args[0])
	},
}

func doDelete(store datastore.Storer, client slack.Client, id string) error {
	a, err := store.GetAnnouncement(id)
	if err != nil {
		return fmt.Errorf("failed to get announcement: %w", err)
	}

	sentMessages, err := store.ListSentMessagesByAnnouncementID(id)
	if err != nil {
		return fmt.Errorf("failed to list sent messages: %w", err)
	}

	for _, sm := range sentMessages {
		if sm.Status == datastore.StatusSent {
			channelID, err := client.GetChannelID(a.ChannelID)
			if err != nil {
				return fmt.Errorf("failed to get channel ID: %w", err)
			}
			if err := client.DeleteMessage(channelID, sm.Timestamp); err != nil {
				return fmt.Errorf("failed to delete message from slack: %w", err)
			}
		}
		if err := store.DeleteSentMessage(sm.ID); err != nil {
			return fmt.Errorf("failed to delete sent message: %w", err)
		}
	}

	if err := store.DeleteAnnouncement(id); err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}

	fmt.Printf("Announcement '%s' and all its sent messages have been deleted.\n", id)

	return nil
}
