/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/andrewhowdencom/announce/cmd/slack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// slackCmd represents the slack command
var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Manage Slack announcements",
	Long:  `Manage Slack announcements.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("slack.app.token") == "" {
			return fmt.Errorf("slack app token is not configured")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(slackCmd)
	slackCmd.AddCommand(slack.ScheduleCmd)
	slackCmd.AddCommand(slack.SendDueCmd)
	slackCmd.AddCommand(slack.WatchCmd)
	slackCmd.AddCommand(slack.DeleteCmd)

	slackCmd.PersistentFlags().String("app-token", "", "Slack app token")
	viper.BindPFlag("slack.app.token", slackCmd.PersistentFlags().Lookup("app-token"))
}
