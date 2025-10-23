package slack

import (
	"fmt"
	"time"

	"github.com/andrewhowdencom/ruf/internal/datastore"
	"github.com/spf13/cobra"
)

// RescheduleCmd represents the reschedule command
var RescheduleCmd = &cobra.Command{
	Use:   "reschedule",
	Short: "Reschedule a message to a Slack channel",
	Long: `Reschedule a message to a Slack channel.

This command reschedules a message to be sent to a Slack channel at a specified time.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the call ID from the flags
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

	// Get the call from the datastore
	call, err := store.GetCall(id)
	if err != nil {
		return fmt.Errorf("failed to get call: %w", err)
	}

	// Update the scheduled time and status
	call.ScheduledAt = scheduledAt
	call.Status = datastore.StatusPending

	// Update the call in the datastore
	if err := store.UpdateCall(call); err != nil {
		return fmt.Errorf("failed to update call: %w", err)
	}

	fmt.Printf("Call with ID %s rescheduled successfully to %s\n", call.ID, call.ScheduledAt.Format(time.RFC3339))
	return nil
}

func init() {
	RescheduleCmd.Flags().String("id", "", "Call ID")
	RescheduleCmd.Flags().String("at", "", "Time to send the message (RFC3339 format)")
}
