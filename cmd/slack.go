/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// slackCmd represents the slack command
var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Post a message to a Slack channel",
	Long: `Reads a message from STDIN and posts it to the configured Slack channel.`,
	Run: func(cmd *cobra.Command, args []string) {
		webhookURL := viper.GetString("slack_webhook_url")
		if webhookURL == "" || webhookURL == "YOUR_WEBHOOK_URL_HERE" {
			fmt.Println("Slack webhook URL is not configured. Please set 'slack_webhook_url' in your config file.")
			os.Exit(1)
		}

		message, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading from STDIN:", err)
			os.Exit(1)
		}

		if len(message) == 0 {
			fmt.Println("No message provided via STDIN.")
			os.Exit(1)
		}

		payload := []byte(fmt.Sprintf(`{"text": "%s"}`, string(message)))
		req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("Error creating request:", err)
			os.Exit(1)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request to Slack:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("Slack API returned a non-200 status code: %d %s\n", resp.StatusCode, string(body))
			os.Exit(1)
		}

		fmt.Println("Message posted to Slack successfully.")
	},
}

func init() {
	// This is intentionally left empty because the command is added in root.go
}