/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/andrewhowdencom/announce/cmd/slack"
	"github.com/spf13/cobra"
)

// slackCmd represents the slack command
var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Manage Slack announcements",
	Long:  `Manage Slack announcements.`,
}

func init() {
	rootCmd.AddCommand(slackCmd)
	slackCmd.AddCommand(slack.PostCmd)
}