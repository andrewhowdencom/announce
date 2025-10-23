package cmd

import (
	"fmt"
	"time"

	"github.com/andrewhowdencom/ruf/internal/clients/slack"
	"github.com/andrewhowdencom/ruf/internal/datastore"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run the worker to send calls",
	Long:  `Run the worker to send calls.`,
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
	calls, err := store.ListCalls()
	if err != nil {
		return fmt.Errorf("failed to list calls: %w", err)
	}

	for _, a := range calls {
		if err := processCall(store, slackClient, a); err != nil {
			fmt.Printf("Error processing call %s: %v\n", a.ID, err)
		}
	}

	return nil
}

func processCall(store datastore.Storer, slackClient slack.Client, a *datastore.Call) error {
	if !((a.Status == datastore.StatusPending || a.Status == datastore.StatusRecurring) && time.Now().After(a.ScheduledAt)) {
		return nil
	}

	fmt.Printf("Sending call %s... ", a.ID)

	timestamp, err := slackClient.PostMessage(a.ChannelID, a.Content)
	sentMessage := &datastore.SentMessage{
		CallID: a.ID,
		Timestamp:      timestamp,
	}

	if err != nil {
		sentMessage.Status = datastore.StatusFailed
		fmt.Printf("failed: %v\n", err)
	} else {
		sentMessage.Status = datastore.StatusSent
		fmt.Println("done")
	}

	if err := store.AddSentMessage(sentMessage); err != nil {
		return fmt.Errorf("failed to add sent message: %w", err)
	}

	if a.Recurring {
		reschedule(a)
		if err := store.UpdateCall(a); err != nil {
			return fmt.Errorf("failed to update call %s: %w", a.ID, err)
		}
	} else {
		a.Status = datastore.StatusProcessed
		if err := store.UpdateCall(a); err != nil {
			return fmt.Errorf("failed to update call %s: %w", a.ID, err)
		}
	}

	return nil
}

func reschedule(a *datastore.Call) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(a.Cron)
	if err != nil {
		// If the cron is invalid, mark as failed
		a.Status = datastore.StatusFailed
		fmt.Printf("failed to parse cron: %v\n", err)
	} else {
		a.ScheduledAt = schedule.Next(time.Now())
	}
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
