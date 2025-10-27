package cmd

import (
	"fmt"
	"log/slog"
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
	slog.Debug("running worker")
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

	slog.Info("starting worker")
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// Run a poll on startup
	if err := runTick(store, slackClient, emailClient, p); err != nil {
		slog.Error("error running tick", "error", err)
	}

	for range ticker.C {
		if err := runTick(store, slackClient, emailClient, p); err != nil {
			slog.Error("error running tick", "error", err)
		}
	}

	return nil
}

func runTick(store datastore.Storer, slackClient slack.Client, emailClient email.Client, p *poller.Poller) error {
	slog.Debug("running tick")
	urls := viper.GetStringSlice("source.urls")
	slog.Debug("polling for calls", "urls", urls)
	calls, err := p.Poll(urls)
	if err != nil {
		return err
	}

	for _, call := range calls {
		if err := processCall(store, slackClient, emailClient, call); err != nil {
			slog.Error("error processing call", "call_id", call.ID, "error", err)
		}
	}

	return nil
}

func processCall(store datastore.Storer, slackClient slack.Client, emailClient email.Client, call *model.Call) error {
	slog.Debug("processing call", "call_id", call.ID)
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
	slog.Debug("calculated effective scheduled time", "call_id", call.ID, "effective_scheduled_at", effectiveScheduledAt)

	// Don't process calls scheduled for the future.
	if now.Before(effectiveScheduledAt) {
		slog.Debug("skipping call scheduled for the future", "call_id", call.ID, "effective_scheduled_at", effectiveScheduledAt)
		return nil
	}

	lookbackPeriod := viper.GetDuration("worker.lookback_period")
	if effectiveScheduledAt.Before(now.Add(-lookbackPeriod)) {
		slog.Warn("skipping call outside lookback period", "call_id", call.ID, "scheduled_at", effectiveScheduledAt)
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
		if len(dest.To) == 0 {
			slog.Warn("skipping call with no address in `to`", "call_id", call.ID)
			continue
		}

		for _, to := range dest.To {
			hasBeenSent, err := store.HasBeenSent(call.ID, effectiveScheduledAt, dest.Type, to)
			if err != nil {
				return fmt.Errorf("failed to check if call has been sent: %w", err)
			}
			if hasBeenSent {
				slog.Debug("skipping call that has already been sent", "call_id", call.ID, "destination", to, "type", dest.Type)
				continue
			}

			switch dest.Type {
			case "slack":
				slog.Info("sending slack message", "call_id", call.ID, "channel", to, "scheduled_at", effectiveScheduledAt)
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
					slog.Error("failed to send slack message", "error", err)
				} else {
					sentMessage.Status = datastore.StatusSent
					slog.Info("sent slack message", "call_id", call.ID, "channel", to, "scheduled_at", effectiveScheduledAt)
				}

				if err := store.AddSentMessage(sentMessage); err != nil {
					return err
				}
			case "email":
				slog.Info("sending email", "call_id", call.ID, "recipient", to, "scheduled_at", effectiveScheduledAt)
				err := emailClient.Send([]string{to}, call.Subject, call.Content)
				sentMessage := &datastore.SentMessage{
					SourceID:    call.ID,
					ScheduledAt: effectiveScheduledAt,
					Destination: to,
					Type:        dest.Type,
				}

				if err != nil {
					sentMessage.Status = datastore.StatusFailed
					slog.Error("failed to send email", "error", err)
				} else {
					sentMessage.Status = datastore.StatusSent
					slog.Info("sent email", "call_id", call.ID, "recipient", to, "scheduled_at", effectiveScheduledAt)
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
