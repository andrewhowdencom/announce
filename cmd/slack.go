/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// slackCmd represents the slack command
var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Post a message to a Slack channel",
	Long: `Reads a message from STDIN and posts it to a Slack channel
configured via a webhook URL.`,
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
	rootCmd.AddCommand(slackCmd)

	slackCmd.Flags().String("slack.webhook-url", "", "Slack webhook URL")
	viper.BindPFlag("slack.webhook_url", slackCmd.Flags().Lookup("slack.webhook-url"))
	viper.SetDefault("slack.webhook_url", "")
}