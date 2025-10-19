package slack

import (
	"fmt"
	"time"

	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/spf13/cobra"
)

// RescheduleCmd represents the reschedule command
var RescheduleCmd = &cobra.Command{
	Use:   "reschedule",
	Short: "Reschedule a message to a Slack channel",
	Long: `Reschedule a message to a Slack channel.

This command reschedules a message to be sent to a Slack channel at a specified time.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the announcement ID from the flags
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		// Get the scheduled time from the flags
		at, err := cmd.Flags().GetString("at")
		if err != nil {
			return err
		}

		// Create a new datastore
		store, err := datastore.NewStore()
		if err != nil {
			return fmt.Errorf("failed to create datastore: %w", err)
		}
		defer store.Close()

		return doReschedule(store, id, at)
	},
}

func doReschedule(store datastore.Storer, id, at string) error {
	scheduledAt, err := time.Parse(time.RFC3339, at)
	if err != nil {
		return fmt.Errorf("failed to parse 'at' flag: %w", err)
	}

	// Get the announcement from the datastore
	announcement, err := store.GetAnnouncement(id)
	if err != nil {
		return fmt.Errorf("failed to get announcement: %w", err)
	}

	// Update the scheduled time and status
	announcement.ScheduledAt = scheduledAt
	announcement.Status = datastore.StatusPending

	// Update the announcement in the datastore
	if err := store.UpdateAnnouncement(announcement); err != nil {
		return fmt.Errorf("failed to update announcement: %w", err)
	}

	fmt.Printf("Announcement with ID %s rescheduled successfully to %s\n", announcement.ID, announcement.ScheduledAt.Format(time.RFC3339))
	return nil
}

func init() {
	RescheduleCmd.Flags().String("id", "", "Announcement ID")
	RescheduleCmd.Flags().String("at", "", "Time to send the message (RFC3339 format)")
}
