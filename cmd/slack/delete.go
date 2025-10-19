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

	if a.Status == datastore.StatusSent {
		// Soft delete: update the status to "deleted" and delete the message from Slack
		channelID, err := client.GetChannelID(a.ChannelID)
		if err != nil {
			return fmt.Errorf("failed to get channel ID: %w", err)
		}
		if err := client.DeleteMessage(channelID, a.Timestamp); err != nil {
			return fmt.Errorf("failed to delete message from slack: %w", err)
		}

		a.Status = datastore.StatusDeleted
		if err := store.UpdateAnnouncement(a); err != nil {
			return fmt.Errorf("failed to update announcement: %w", err)
		}
		fmt.Printf("Announcement '%s' soft deleted and message removed from Slack.\n", id)
	} else {
		// Hard delete: remove the announcement from the datastore
		if err := store.DeleteAnnouncement(id); err != nil {
			return fmt.Errorf("failed to delete announcement: %w", err)
		}
		fmt.Printf("Announcement '%s' hard deleted.\n", id)
	}

	return nil
}
