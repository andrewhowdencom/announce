package model

import "time"

// Destination represents a destination to send a call to.
type Destination struct {
	Type      string `json:"type"`
	ChannelID string `json:"channel_id"`
}

// Call represents a message to be sent to a destination.
type Call struct {
	ID           string        `json:"id"`
	Content      string        `json:"content"`
	Destinations []Destination `json:"destinations"`
	ScheduledAt  time.Time     `json:"scheduled_at"`
	Cron         string        `json:"cron,omitempty"`
	Recurring    bool          `json:"recurring,omitempty"`
}
