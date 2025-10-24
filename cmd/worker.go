package cmd

import (
	"fmt"
	"time"

	"github.com/andrewhowdencom/ruf/internal/clients/slack"
	"github.com/andrewhowdencom/ruf/internal/datastore"
	"github.com/andrewhowdencom/ruf/internal/model"
	"github.com/andrewhowdencom/ruf/internal/sourcer"
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
		return runWorker()
	},
}

func buildSourcer() sourcer.Sourcer {
	fetcher := sourcer.NewCompositeFetcher()
	fetcher.AddFetcher("http", sourcer.NewHTTPFetcher())
	fetcher.AddFetcher("https", sourcer.NewHTTPFetcher())
	fetcher.AddFetcher("file", sourcer.NewFileFetcher())
	parser := sourcer.NewYAMLParser()
	return sourcer.NewSourcer(fetcher, parser)
}

func runWorker() error {
	store, err := datastore.NewStore()
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}
	defer store.Close()

	slackToken := viper.GetString("slack.app.token")
	slackClient := slack.NewClient(slackToken)

	s := buildSourcer()

	fmt.Println("Starting worker...")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := runTick(store, slackClient, s); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	return nil
}

func runTick(store datastore.Storer, slackClient slack.Client, s sourcer.Sourcer) error {
	urls := viper.GetStringSlice("source.urls")
	var allCalls []*model.Call

	for _, url := range urls {
		calls, err := s.Source(url)
		if err != nil {
			fmt.Printf("Error sourcing from %s: %v\n", url, err)
			continue
		}
		allCalls = append(allCalls, calls...)
	}

	for _, call := range allCalls {
		if err := processCall(store, slackClient, call); err != nil {
			fmt.Printf("Error processing call %s: %v\n", call.ID, err)
		}
	}

	return nil
}

func processCall(store datastore.Storer, slackClient slack.Client, call *model.Call) error {
	now := time.Now()
	var effectiveScheduledAt time.Time

	if call.Cron != "" {
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(call.Cron)
		if err != nil {
			return fmt.Errorf("failed to parse cron: %w", err)
		}
		// Find the last scheduled time before or at `now`.
		effectiveScheduledAt = schedule.Next(now.Add(-2 * time.Minute)).Truncate(time.Minute)

	} else {
		effectiveScheduledAt = call.ScheduledAt
	}

	// Don't process calls scheduled for the future.
	if now.Before(effectiveScheduledAt) {
		return nil
	}

	hasBeenSent, err := store.HasBeenSent(call.ID, effectiveScheduledAt)
	if err != nil {
		return fmt.Errorf("failed to check if call has been sent: %w", err)
	}
	if hasBeenSent {
		return nil
	}

	for _, dest := range call.Destinations {
		switch dest.Type {
		case "slack":
			fmt.Printf("Sending call %s to Slack channel %s for %v... ", call.ID, dest.ChannelID, effectiveScheduledAt)

			timestamp, err := slackClient.PostMessage(dest.ChannelID, call.Content)
			sentMessage := &datastore.SentMessage{
				SourceID:    call.ID,
				ScheduledAt: effectiveScheduledAt,
				Timestamp:   timestamp,
				ChannelID:   dest.ChannelID,
			}

			if err != nil {
				sentMessage.Status = datastore.StatusFailed
				fmt.Printf("failed: %v\n", err)
			} else {
				sentMessage.Status = datastore.StatusSent
				fmt.Println("done")
			}

			if err := store.AddSentMessage(sentMessage); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported destination type: %s", dest.Type)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
