package datastore

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"go.etcd.io/bbolt"
)

// Err* are common errors returned by the datastore.
var (
	ErrNotFound            = errors.New("not found")
	ErrDBOperationFailed   = errors.New("db operation failed")
	ErrSerializationFailed = errors.New("serialization failed")
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
	ID           string    `json:"id"`
	SourceID     string    `json:"source_id"`
	ScheduledAt  time.Time `json:"scheduled_at"`
	Timestamp    string    `json:"timestamp,omitempty"` // Slack timestamp
	Destination  string    `json:"destination"`
	Type         string    `json:"type"`
	Status       Status    `json:"status"`
	CampaignName string    `json:"campaign_name"`
}

// Storer is an interface that defines the methods for interacting with the datastore.
type Storer interface {
	AddSentMessage(campaignID, callID string, sm *SentMessage) error
	HasBeenSent(campaignID, callID, destType, destination string) (bool, error)
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
		return nil, fmt.Errorf("%w: failed to get db path: %w", ErrDBOperationFailed, err)
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
		return nil, fmt.Errorf("%w: failed to open db: %w", ErrDBOperationFailed, err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(sentMessagesBucket)
		if err != nil {
			return fmt.Errorf("%w: failed to create bucket: %w", ErrDBOperationFailed, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// AddSentMessage adds a new sent message to the store.
func (s *Store) AddSentMessage(campaignID, callID string, sm *SentMessage) error {
	err := s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(sentMessagesBucket)
		sm.ID = s.generateID(campaignID, callID, sm.Type, sm.Destination)

		buf, err := json.Marshal(sm)
		if err != nil {
			return fmt.Errorf("%w: failed to marshal sent message: %w", ErrSerializationFailed, err)
		}

		if err := b.Put([]byte(sm.ID), buf); err != nil {
			return fmt.Errorf("%w: failed to put sent message: %w", ErrDBOperationFailed, err)
		}
		return nil
	})
	return err
}

// HasBeenSent checks if a message with the given sourceID and scheduledAt time has a 'sent' or 'deleted' status.
// It returns false for messages that have a 'failed' status, or do not exist.
func (s *Store) HasBeenSent(campaignID, callID, destType, destination string) (bool, error) {
	var sent bool
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(sentMessagesBucket)
		id := s.generateID(campaignID, callID, destType, destination)
		v := b.Get([]byte(id))
		if v != nil {
			var sm SentMessage
			if err := json.Unmarshal(v, &sm); err != nil {
				return fmt.Errorf("%w: failed to unmarshal sent message: %w", ErrSerializationFailed, err)
			}
			if sm.Status == StatusSent || sm.Status == StatusDeleted {
				sent = true
			}
		}
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("%w: failed to check if message has been sent: %w", ErrDBOperationFailed, err)
	}
	return sent, nil
}

func (s *Store) generateID(campaignID, callID, destType, destination string) string {
	parts := []string{
		campaignID,
		callID,
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
		err := b.ForEach(func(k, v []byte) error {
			var sm SentMessage
			if err := json.Unmarshal(v, &sm); err != nil {
				return fmt.Errorf("%w: failed to unmarshal sent message: %w", ErrSerializationFailed, err)
			}
			sentMessages = append(sentMessages, &sm)
			return nil
		})
		if err != nil {
			return fmt.Errorf("%w: failed to iterate over sent messages: %w", ErrDBOperationFailed, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
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
			return fmt.Errorf("%w: message with id '%s'", ErrNotFound, id)
		}
		if err := json.Unmarshal(v, &sm); err != nil {
			return fmt.Errorf("%w: failed to unmarshal sent message: %w", ErrSerializationFailed, err)
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
			return fmt.Errorf("%w: message with id '%s'", ErrNotFound, id)
		}

		var sm SentMessage
		if err := json.Unmarshal(v, &sm); err != nil {
			return fmt.Errorf("%w: failed to unmarshal sent message: %w", ErrSerializationFailed, err)
		}

		sm.Status = StatusDeleted

		buf, err := json.Marshal(sm)
		if err != nil {
			return fmt.Errorf("%w: failed to marshal sent message: %w", ErrSerializationFailed, err)
		}

		if err := b.Put([]byte(id), buf); err != nil {
			return fmt.Errorf("%w: failed to put sent message: %w", ErrDBOperationFailed, err)
		}
		return nil
	})
}
