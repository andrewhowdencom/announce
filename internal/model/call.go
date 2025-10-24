package model

import "time"

// Destination represents a destination to send a call to.
type Destination struct {
	Type      string `json:"type" yaml:"type"`
	ChannelID string `json:"channel_id" yaml:"channel_id"`
}

// Call represents a message to be sent to a destination.
type Call struct {
	ID           string        `json:"id" yaml:"id"`
	Content      string        `json:"content" yaml:"content"`
	Destinations []Destination `json:"destinations" yaml:"destinations"`
	ScheduledAt  time.Time     `json:"scheduled_at" yaml:"scheduled_at"`
	Cron         string        `json:"cron,omitempty" yaml:"cron,omitempty"`
	Recurring    bool          `json:"recurring,omitempty" yaml:"recurring,omitempty"`
}
