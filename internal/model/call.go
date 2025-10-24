package model

import "time"

// Destination represents a destination to send a call to.
type Destination struct {
	Type      string `json:"type"`
	ChannelID string `json:"channel_id"`
}

// Email represents an email to be sent.
type Email struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

// Call represents a message to be sent to a destination.
type Call struct {
	ID           string        `json:"id"`
	Content      string        `json:"content"`
	Destinations []Destination `json:"destinations"`
	Email        *Email        `json:"email,omitempty"`
	ScheduledAt  time.Time     `json:"scheduled_at"`
	Cron         string        `json:"cron,omitempty"`
	Recurring    bool          `json:"recurring,omitempty"`
}
