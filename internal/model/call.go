package model

import "time"

// Destination represents a destination to send a call to.
type Destination struct {
	Type      string `json:"type" yaml:"type"`
	ChannelID string `json:"channel_id" yaml:"channel_id"`
}

// Email represents an email to be sent.
type Email struct {
	To      []string `json:"to" yaml:"to"`
	Subject string   `json:"subject" yaml:"subject"`
	Body    string   `json:"body" yaml:"body"`
}

// Call represents a message to be sent to a destination.
type Call struct {
	ID           string        `json:"id" yaml:"id"`
	Content      string        `json:"content" yaml:"content"`
	Destinations []Destination `json:"destinations" yaml:"destinations"`
	ScheduledAt  time.Time     `json:"scheduled_at" yaml:"scheduled_at"`
	Cron         string        `json:"cron,omitempty" yaml:"cron,omitempty"`
	Recurring    bool          `json:"recurring,omitempty" yaml:"recurring,omitempty"`
	Email        *Email        `json:"email,omitempty" yaml:"email,omitempty"`
}
