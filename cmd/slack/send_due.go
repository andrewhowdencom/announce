package slack

import (
	"fmt"
	"time"

	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SendDueCmd represents the send-due command
var SendDueCmd = &cobra.Command{
	Use:   "send-due",
	Short: "Send all due announcements",
	Long:  `Send all due announcements.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a new datastore
		store, err := datastore.NewStore()
		if err != nil {
			return fmt.Errorf("failed to create datastore: %w", err)
		}
		defer store.Close()

		// List all announcements
		announcements, err := store.ListAnnouncements()
		if err != nil {
			return fmt.Errorf("failed to list announcements: %w", err)
		}

		// Create a new Slack client
		client := NewSlackClient(viper.GetString("slack.app.token"))

		// Iterate over the announcements and send the due ones
		for _, a := range announcements {
			if a.Status == datastore.StatusPending && time.Now().After(a.ScheduledAt) {
				fmt.Printf("Sending announcement %s... ", a.ID)
				timestamp, err := client.PostMessage(a.ChannelID, a.Content)
				if err != nil {
					a.Status = datastore.StatusFailed
					fmt.Printf("failed: %v\n", err)
				} else {
					a.Status = datastore.StatusSent
					a.Timestamp = timestamp
					fmt.Println("done")
				}
				if err := store.UpdateAnnouncement(a); err != nil {
					return fmt.Errorf("failed to update announcement %s: %w", a.ID, err)
				}
			}
		}

		return nil
	},
}
