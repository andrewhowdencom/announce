package slack

import (
	"fmt"
	"time"

	"github.com/andrewhowdencom/announce/internal/clients/slack"
	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SendDueCmd represents the send-due command
var SendDueCmd = &cobra.Command{
	Use:   "send-due",
	Short: "Send all due announcements",
	Long:  `Send all due announcements.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a new datastore
		store, err := datastore.NewStore()
		if err != nil {
			return fmt.Errorf("failed to create datastore: %w", err)
		}
		defer store.Close()

		// List all announcements
		announcements, err := store.ListAnnouncements()
		if err != nil {
			return fmt.Errorf("failed to list announcements: %w", err)
		}

		// Create a new Slack client
		client := slack.NewClient(viper.GetString("slack.app.token"))

		// Iterate over the announcements and send the due ones
		for _, a := range announcements {
			if (a.Status == datastore.StatusPending || a.Status == datastore.StatusRecurring) && time.Now().After(a.ScheduledAt) {
				fmt.Printf("Sending announcement %s... ", a.ID) // Use raw string literal for the format string

				channelID, err := client.GetChannelID(a.ChannelID)
				if err != nil {
					a.Status = datastore.StatusFailed
					fmt.Printf("failed to get channel ID: %v\n", err)
				} else {
					timestamp, err := client.PostMessage(channelID, a.Content)
					if err != nil {
						a.Status = datastore.StatusFailed
						fmt.Printf("failed: %v\n", err)
					} else {
						if a.Recurring {
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
						} else {
							a.Status = datastore.StatusSent
							a.Timestamp = timestamp
							fmt.Println("done")
						}
					}
				}

				if err := store.UpdateAnnouncement(a); err != nil {
					return fmt.Errorf("failed to update announcement %s: %w", a.ID, err)
				}
			}
		}

		return nil
	},
}