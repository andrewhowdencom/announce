package datastore

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"go.etcd.io/bbolt"
)

var sentMessagesBucket = []byte("sent_messages")

// Status represents the status of a call.
type Status string

const (
	// StatusSent means the call has been successfully sent.
	StatusSent Status = "sent"
	// StatusFailed means the call failed to send.
	StatusFailed Status = "failed"
	// StatusDeleted means the call has been deleted.
	StatusDeleted Status = "deleted"
)

// SentMessage represents a message that has been sent.
type SentMessage struct {
	ID          string    `json:"id"`
	SourceID    string    `json:"source_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Timestamp   string    `json:"timestamp,omitempty"` // Slack timestamp
	Destination string    `json:"destination"`
	Type        string    `json:"type"`
	Status      Status    `json:"status"`
}

// Storer is an interface that defines the methods for interacting with the datastore.
type Storer interface {
	AddSentMessage(sm *SentMessage) error
	HasBeenSent(sourceID string, scheduledAt time.Time, destType, destination string) (bool, error)
	ListSentMessages() ([]*SentMessage, error)
	GetSentMessage(id string) (*SentMessage, error)
	DeleteSentMessage(id string) error
	Close() error
}

// Store manages the persistence of calls.
type Store struct {
	db *bbolt.DB
}

// NewStore creates a new Store and initializes the database.
func NewStore() (Storer, error) {
	dbPath, err := xdg.DataFile("ruf/ruf.db")
	if err != nil {
		return nil, fmt.Errorf("failed to get db path: %w", err)
	}

	return newStore(dbPath)
}

// NewTestStore creates a new Store for testing purposes.
func NewTestStore(dbPath string) (Storer, error) {
	return newStore(dbPath)
}

func newStore(dbPath string) (Storer, error) {
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(sentMessagesBucket)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// AddSentMessage adds a new sent message to the store.
func (s *Store) AddSentMessage(sm *SentMessage) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(sentMessagesBucket)
		sm.ID = s.generateID(sm.SourceID, sm.ScheduledAt, sm.Type, sm.Destination)

		buf, err := json.Marshal(sm)
		if err != nil {
			return fmt.Errorf("failed to marshal sent message: %w", err)
		}
		return b.Put([]byte(sm.ID), buf)
	})
}

// HasBeenSent checks if a message with the given sourceID and scheduledAt time has been sent.
func (s *Store) HasBeenSent(sourceID string, scheduledAt time.Time, destType, destination string) (bool, error) {
	var sent bool
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(sentMessagesBucket)
		id := s.generateID(sourceID, scheduledAt, destType, destination)
		v := b.Get([]byte(id))
		if v != nil {
			var sm SentMessage
			if err := json.Unmarshal(v, &sm); err != nil {
				return fmt.Errorf("failed to unmarshal sent message: %w", err)
			}
			if sm.Status != StatusDeleted {
				sent = true
			}
		}
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("failed to check if message has been sent: %w", err)
	}
	return sent, nil
}

func (s *Store) generateID(sourceID string, scheduledAt time.Time, destType, destination string) string {
	parts := []string{
		sourceID,
		scheduledAt.Format(time.RFC3339Nano),
		destType,
		destination,
	}
	return strings.Join(parts, "@")
}

// ListSentMessages retrieves all sent messages from the store.
func (s *Store) ListSentMessages() ([]*SentMessage, error) {
	var sentMessages []*SentMessage
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(sentMessagesBucket)
		return b.ForEach(func(k, v []byte) error {
			var sm SentMessage
			if err := json.Unmarshal(v, &sm); err != nil {
				return fmt.Errorf("failed to unmarshal sent message: %w", err)
			}
			sentMessages = append(sentMessages, &sm)
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list sent messages: %w", err)
	}
	return sentMessages, nil
}

// GetSentMessage retrieves a single sent message from the store.
func (s *Store) GetSentMessage(id string) (*SentMessage, error) {
	var sm SentMessage
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(sentMessagesBucket)
		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("message with id '%s' not found", id)
		}
		if err := json.Unmarshal(v, &sm); err != nil {
			return fmt.Errorf("failed to unmarshal sent message: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &sm, nil
}

// DeleteSentMessage removes a sent message from the store.
func (s *Store) DeleteSentMessage(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(sentMessagesBucket)
		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("message with id '%s' not found", id)
		}

		var sm SentMessage
		if err := json.Unmarshal(v, &sm); err != nil {
			return fmt.Errorf("failed to unmarshal sent message: %w", err)
		}

		sm.Status = StatusDeleted

		buf, err := json.Marshal(sm)
		if err != nil {
			return fmt.Errorf("failed to marshal sent message: %w", err)
		}
		return b.Put([]byte(id), buf)
	})
}
