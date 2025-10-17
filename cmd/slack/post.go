/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package slack

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PostCmd represents the post command
var PostCmd = &cobra.Command{
	Use:   "post",
	Short: "Post a message to a Slack channel",
	Long: `Post a message to a Slack channel.

This command reads a message from STDIN and posts it to a Slack channel
configured via a webhook URL.

Example:
  echo "Hello, world!" | announce slack post --webhook-url <your-webhook-url>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read the message from STDIN
		message, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read message from STDIN: %w", err)
		}

		// Get the webhook URL from the configuration
		webhookURL := viper.GetString("slack.webhook_url")
		if webhookURL == "" {
			return fmt.Errorf("slack webhook URL is not configured")
		}

		// Send the message to the webhook URL
		payload := bytes.NewBufferString(fmt.Sprintf(`{"text": "%s"}`, string(message)))
		resp, err := http.Post(webhookURL, "application/json", payload)
		if err != nil {
			return fmt.Errorf("failed to send message to Slack: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to send message to Slack: received status code %d", resp.StatusCode)
		}

		fmt.Println("Message sent to Slack successfully")
		return nil
	},
}

func init() {
	PostCmd.Flags().String("webhook-url", "", "Slack webhook URL")
	viper.BindPFlag("slack.webhook_url", PostCmd.Flags().Lookup("webhook-url"))
	viper.SetDefault("slack.webhook_url", "")
}