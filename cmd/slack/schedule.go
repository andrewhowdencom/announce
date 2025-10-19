package slack

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/spf13/cobra"
)

// ScheduleCmd represents the schedule command
var ScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule a message to a Slack channel",
	Long: `Schedule a message to a Slack channel.

This command reads a message from STDIN and schedules it to be sent to a Slack channel at a specified time.

Example:
  echo "Hello, world!" | announce slack schedule --channel C1234567890 --at 2025-10-26T19:00:00Z`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read the message from STDIN
		message, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read message from STDIN: %w", err)
		}

		// Get the channel ID from the flags
		channelID, err := cmd.Flags().GetString("channel")
		if err != nil {
			return err
		}

		// Get the scheduled time from the flags
		at, err := cmd.Flags().GetString("at")
		if err != nil {
			return err
		}
		scheduledAt, err := time.Parse(time.RFC3339, at)
		if err != nil {
			return fmt.Errorf("failed to parse 'at' flag: %w", err)
		}

		// Create a new datastore
		store, err := datastore.NewStore()
		if err != nil {
			return fmt.Errorf("failed to create datastore: %w", err)
		}
		defer store.Close()

		// Create a new announcement
		announcement := &datastore.Announcement{
			Content:     string(message),
			ChannelID:   channelID,
			ScheduledAt: scheduledAt,
			Status:      datastore.StatusPending,
		}

		// Add the announcement to the datastore
		if err := store.AddAnnouncement(announcement); err != nil {
			return fmt.Errorf("failed to add announcement: %w", err)
		}

		fmt.Printf("Announcement scheduled successfully with ID: %s\n", announcement.ID)
		return nil
	},
}

func init() {
	ScheduleCmd.Flags().String("channel", "", "Slack channel ID")
	ScheduleCmd.Flags().String("at", "", "Time to send the message (RFC3339 format)")
}
