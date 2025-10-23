package slack

import (
	"fmt"

	"github.com/andrewhowdencom/ruf/internal/clients/slack"
	"github.com/andrewhowdencom/ruf/internal/datastore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a call",
	Long:  `Delete a call.`,
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
	a, err := store.GetCall(id)
	if err != nil {
		return fmt.Errorf("failed to get call: %w", err)
	}

	sentMessages, err := store.ListSentMessagesByCallID(id)
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

	if err := store.DeleteCall(id); err != nil {
		return fmt.Errorf("failed to delete call: %w", err)
	}

	fmt.Printf("Call '%s' and all its sent messages have been deleted.\n", id)

	return nil
}
