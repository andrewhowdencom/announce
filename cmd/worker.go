package cmd

import (
	"fmt"
	"time"

	"github.com/andrewhowdencom/ruf/internal/clients/email"
	"github.com/andrewhowdencom/ruf/internal/clients/slack"
	"github.com/andrewhowdencom/ruf/internal/datastore"
	"github.com/andrewhowdencom/ruf/internal/model"
	"github.com/andrewhowdencom/ruf/internal/poller"
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

	emailClient := email.NewClient(
		viper.GetString("email.host"),
		viper.GetInt("email.port"),
		viper.GetString("email.username"),
		viper.GetString("email.password"),
		viper.GetString("email.from"),
	)

	s := buildSourcer()
	pollInterval := viper.GetDuration("worker.interval")
	if pollInterval == 0 {
		pollInterval = 1 * time.Minute
	}
	p := poller.New(s, pollInterval)

	fmt.Println("Starting worker...")
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// Run a poll on startup
	if err := runTick(store, slackClient, emailClient, p); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	for range ticker.C {
		if err := runTick(store, slackClient, emailClient, p); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	return nil
}

func runTick(store datastore.Storer, slackClient slack.Client, emailClient email.Client, p *poller.Poller) error {
	urls := viper.GetStringSlice("source.urls")
	calls, err := p.Poll(urls)
	if err != nil {
		return err
	}

	for _, call := range calls {
		if err := processCall(store, slackClient, emailClient, call); err != nil {
			fmt.Printf("Error processing call %s: %v\n", call.ID, err)
		}
	}

	return nil
}

func processCall(store datastore.Storer, slackClient slack.Client, emailClient email.Client, call *model.Call) error {
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

	lookbackPeriod := viper.GetDuration("worker.lookback_period")
	if effectiveScheduledAt.Before(now.Add(-lookbackPeriod)) {
		fmt.Printf("Skipping call %s scheduled at %s because it is outside the lookback period\n", call.ID, effectiveScheduledAt)
		for _, dest := range call.Destinations {
			for _, to := range dest.To {
				err := store.AddSentMessage(&datastore.SentMessage{
					SourceID:    call.ID,
					ScheduledAt: effectiveScheduledAt,
					Status:      datastore.StatusFailed,
					Type:        dest.Type,
					Destination: to,
				})
				if err != nil {
					return fmt.Errorf("failed to add sent message: %w", err)
				}
			}
		}
		return nil
	}

	for _, dest := range call.Destinations {
		for _, to := range dest.To {
			hasBeenSent, err := store.HasBeenSent(call.ID, effectiveScheduledAt, dest.Type, to)
			if err != nil {
				return fmt.Errorf("failed to check if call has been sent: %w", err)
			}
			if hasBeenSent {
				continue
			}

			switch dest.Type {
			case "slack":
				fmt.Printf("Sending call %s to Slack channel %s for %v... ", call.ID, to, effectiveScheduledAt)

				timestamp, err := slackClient.PostMessage(to, call.Subject, call.Content)
				sentMessage := &datastore.SentMessage{
					SourceID:    call.ID,
					ScheduledAt: effectiveScheduledAt,
					Timestamp:   timestamp,
					Destination: to,
					Type:        dest.Type,
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
			case "email":
				fmt.Printf("Sending call %s to email %s for %v... ", call.ID, to, effectiveScheduledAt)

				err := emailClient.Send([]string{to}, call.Subject, call.Content)
				sentMessage := &datastore.SentMessage{
					SourceID:    call.ID,
					ScheduledAt: effectiveScheduledAt,
					Destination: to,
					Type:        dest.Type,
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
	}

	return nil
}

func init() {
	rootCmd.AddCommand(workerCmd)
	viper.SetDefault("worker.interval", "1m")
	viper.SetDefault("worker.lookback_period", "24h")
}
