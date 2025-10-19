package cmd

import (
	"fmt"
	"time"

	"github.com/andrewhowdencom/announce/internal/clients/slack"
	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run the worker to send announcements",
	Long:  `Run the worker to send announcements.`, // Corrected: Removed unnecessary escaping of backticks
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := datastore.NewStore()
		if err != nil {
			return fmt.Errorf("failed to create store: %w", err)
		}
		defer store.Close()

		slackToken := viper.GetString("slack.app.token")
		slackClient := slack.NewClient(slackToken)

		fmt.Println("Starting worker...")
		for {
			if err := runWorker(store, slackClient); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			time.Sleep(1 * time.Minute)
		}
	},
}

func runWorker(store datastore.Storer, slackClient slack.Client) error {
	announcements, err := store.ListAnnouncements()
	if err != nil {
		return fmt.Errorf("failed to list announcements: %w", err)
	}

	for _, a := range announcements {
		if err := processAnnouncement(store, slackClient, a); err != nil {
			fmt.Printf("Error processing announcement %s: %v\n", a.ID, err)
		}
	}

	return nil
}

func processAnnouncement(store datastore.Storer, slackClient slack.Client, a *datastore.Announcement) error {
	if !shouldSend(a) {
		return nil
	}

	fmt.Printf("Sending announcement %s... ", a.ID)

	timestamp, err := slackClient.PostMessage(a.ChannelID, a.Content)
	if err != nil {
		a.Status = datastore.StatusFailed
		fmt.Printf("failed: %v\n", err)
	} else {
		if a.Recurring {
			reschedule(a)
		} else {
			a.Status = datastore.StatusSent
			a.Timestamp = timestamp
			fmt.Println("done")
		}
	}

	if err := store.UpdateAnnouncement(a); err != nil {
		return fmt.Errorf("failed to update announcement %s: %w", a.ID, err)
	}

	return nil
}

func shouldSend(a *datastore.Announcement) bool {
	return (a.Status == datastore.StatusPending || a.Status == datastore.StatusRecurring) && time.Now().After(a.ScheduledAt)
}

func reschedule(a *datastore.Announcement) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(a.Cron)
	if err != nil {
		// If the cron is invalid, mark as failed
		a.Status = datastore.StatusFailed
		fmt.Printf("failed to parse cron: %v\n", err)
	} else {
		a.ScheduledAt = schedule.Next(time.Now())
		a.Status = datastore.StatusRecurring
		fmt.Println("done, rescheduled")
	}
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
