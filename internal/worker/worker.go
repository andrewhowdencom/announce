package worker

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/andrewhowdencom/ruf/internal/clients/email"
	"github.com/andrewhowdencom/ruf/internal/clients/slack"
	"github.com/andrewhowdencom/ruf/internal/datastore"
	"github.com/andrewhowdencom/ruf/internal/model"
	"github.com/andrewhowdencom/ruf/internal/poller"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

// Worker is responsible for polling for calls and sending them.
type Worker struct {
	store       datastore.Storer
	slackClient slack.Client
	emailClient email.Client
	poller      *poller.Poller
	interval    time.Duration
}

// New creates a new worker.
func New(store datastore.Storer, slackClient slack.Client, emailClient email.Client, poller *poller.Poller, interval time.Duration) *Worker {
	return &Worker{
		store:       store,
		slackClient: slackClient,
		emailClient: emailClient,
		poller:      poller,
		interval:    interval,
	}
}

// Run starts the worker.
func (w *Worker) Run() error {
	slog.Info("starting worker")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run a poll on startup
	if err := w.RunTick(); err != nil {
		slog.Error("error running tick", "error", err)
	}

	for range ticker.C {
		if err := w.RunTick(); err != nil {
			slog.Error("error running tick", "error", err)
		}
	}
	return nil
}

// RunTick performs a single poll for calls and sends them.
func (w *Worker) RunTick() error {
	slog.Debug("running tick")
	urls := viper.GetStringSlice("source.urls")
	slog.Debug("polling for calls", "urls", urls)
	calls, err := w.poller.Poll(urls)
	if err != nil {
		return err
	}

	for _, call := range calls {
		if err := w.processCall(call); err != nil {
			slog.Error("error processing call", "call_id", call.ID, "error", err)
		}
	}

	return nil
}

func (w *Worker) processCall(call *model.Call) error {
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
				err := w.store.AddSentMessage(&datastore.SentMessage{
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
			hasBeenSent, err := w.store.HasBeenSent(call.ID, effectiveScheduledAt, dest.Type, to)
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
				timestamp, err := w.slackClient.PostMessage(to, call.Subject, call.Content)
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

				if err := w.store.AddSentMessage(sentMessage); err != nil {
					return err
				}
			case "email":
				slog.Info("sending email", "call_id", call.ID, "recipient", to, "scheduled_at", effectiveScheduledAt)
				err := w.emailClient.Send([]string{to}, call.Subject, call.Content)
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

				if err := w.store.AddSentMessage(sentMessage); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported destination type: %s", dest.Type)
			}
		}
	}

	return nil
}
