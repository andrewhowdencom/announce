/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package slack

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PostCmd represents the post command
var PostCmd = &cobra.Command{
	Use:   "post",
	Short: "Post a message to a Slack channel",
	Long: `Post a message to a Slack channel.

This command reads a message from STDIN and posts it to a Slack channel.

Example:
  echo "Hello, world!" | announce slack post --channel C1234567890`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read the message from STDIN
		message, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read message from STDIN: %w", err)
		}

		// Get the bot token from the configuration
		botToken := viper.GetString("slack.bot_token")
		if botToken == "" {
			return fmt.Errorf("slack bot token is not configured")
		}

		// Get the channel ID from the flags
		channelID, err := cmd.Flags().GetString("channel")
		if err != nil {
			return err
		}

		// Create a new Slack client
		client := NewSlackClient(botToken)

		// Post the message to the Slack channel
		timestamp, err := client.PostMessage(channelID, string(message))
		if err != nil {
			return fmt.Errorf("failed to send message to Slack: %w", err)
		}

		fmt.Printf("Message sent to Slack successfully with timestamp: %s\n", timestamp)
		return nil
	},
}

func init() {
	PostCmd.Flags().String("channel", "", "Slack channel ID")
	PostCmd.Flags().String("bot-token", "", "Slack bot token")
	viper.BindPFlag("slack.bot_token", PostCmd.Flags().Lookup("bot-token"))
}