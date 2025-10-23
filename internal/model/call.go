package model

import "time"

// Call represents a message to be sent to a destination.
type Call struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	ChannelID   string    `json:"channel_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Cron        string    `json:"cron,omitempty"`
	Recurring   bool      `json:"recurring,omitempty"`
}
