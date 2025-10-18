package slack

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a message from a Slack channel",
	Long: `Delete a message from a Slack channel.

This command deletes a message from a Slack channel given a channel ID and a message timestamp.

Example:
  announce slack delete --channel C1234567890 --timestamp 1234567890.123456`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the channel ID from the flags
		channelID, err := cmd.Flags().GetString("channel")
		if err != nil {
			return err
		}

		// Get the timestamp from the flags
		timestamp, err := cmd.Flags().GetString("timestamp")
		if err != nil {
			return err
		}

		// Create a new Slack client
		client := NewSlackClient(viper.GetString("slack.app.token"))

		// Delete the message from the Slack channel
		if err := client.DeleteMessage(channelID, timestamp); err != nil {
			return fmt.Errorf("failed to delete message from Slack: %w", err)
		}

		fmt.Println("Message deleted from Slack successfully")
		return nil
	},
}

func init() {
	DeleteCmd.Flags().String("channel", "", "Slack channel ID")
	DeleteCmd.Flags().String("timestamp", "", "Timestamp of the message to delete")
}