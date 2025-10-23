package datastore

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adrg/xdg"
	"go.etcd.io/bbolt"
)

// The name of the bucket where calls will be stored.
var callsBucket = []byte("calls")
var sentMessagesBucket = []byte("sent_messages")

// Status represents the status of an call.
type Status string

const (
	// StatusPending means the call has been scheduled but not yet sent.
	StatusPending Status = "pending"
	// StatusSent means the call has been successfully sent.
	StatusSent Status = "sent"
	// StatusFailed means the call failed to send.
	StatusFailed Status = "failed"
	// StatusDeleted means the call has been deleted.
	StatusDeleted Status = "deleted"
	// StatusRecurring means the call is recurring.
	StatusRecurring Status = "recurring"
	// StatusProcessed means the call has been processed.
	StatusProcessed Status = "processed"
)

// Call represents a message to be sent to a destination.
type Call struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	ChannelID   string    `json:"channel_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      Status    `json:"status"`
	Cron        string    `json:"cron,omitempty"`
	Recurring   bool      `json:"recurring,omitempty"`
}

// SentMessage represents a message that has been sent.
type SentMessage struct {
	ID             string    `json:"id"`
	CallID string    `json:"call_id"`
	Timestamp      string    `json:"timestamp"`
	Status         Status    `json:"status"`
}

// Storer is an interface that defines the methods for interacting with the datastore.
type Storer interface {
	AddCall(a *Call) error
	GetCall(id string) (*Call, error)
	ListCalls() ([]*Call, error)
	UpdateCall(a *Call) error
	DeleteCall(id string) error

	// AddSentMessage adds a new sent message to the datastore.
	AddSentMessage(sm *SentMessage) error
	// ListSentMessages lists all sent messages from the datastore.
	ListSentMessages() ([]*SentMessage, error)
	// ListSentMessagesByCallID lists all sent messages for a given call ID.
	ListSentMessagesByCallID(callID string) ([]*SentMessage, error)
	// DeleteSentMessage deletes a sent message from the datastore.
	DeleteSentMessage(id string) error

	Close() error
}

// Store manages the persistence of calls.
type Store struct {
	db *bbolt.DB
}

// NewStore creates a new Store and initializes the database.
func NewStore() (Storer, error) {
	dbPath, err := xdg.DataFile("announce/announce.db")
	if err != nil {
		return nil, fmt.Errorf("failed to get db path: %w", err)
	}

	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(callsBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(sentMessagesBucket)
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

// AddCall adds a new call to the store.
func (s *Store) AddCall(a *Call) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(callsBucket)
		id, _ := b.NextSequence()
		a.ID = fmt.Sprintf("%d", id)

		buf, err := json.Marshal(a)
		if err != nil {
			return fmt.Errorf("failed to marshal call: %w", err)
		}
		return b.Put([]byte(a.ID), buf)
	})
}

// ListCalls retrieves all calls from the store.
func (s *Store) ListCalls() ([]*Call, error) {
	var calls []*Call
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(callsBucket)
		return b.ForEach(func(k, v []byte) error {
			var a Call
			if err := json.Unmarshal(v, &a); err != nil {
				return fmt.Errorf("failed to unmarshal call: %w", err)
			}
			calls = append(calls, &a)
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list calls: %w", err)
	}
	return calls, nil
}

// UpdateCall updates an existing call in the store.
func (s *Store) UpdateCall(a *Call) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(callsBucket)
		buf, err := json.Marshal(a)
		if err != nil {
			return fmt.Errorf("failed to marshal call: %w", err)
		}
		return b.Put([]byte(a.ID), buf)
	})
}

// GetCall retrieves a single call from the store.
func (s *Store) GetCall(id string) (*Call, error) {
	var a *Call
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(callsBucket)
		v := b.Get([]byte(id))
		if v == nil {
			return ErrNotFound
		}
		if err := json.Unmarshal(v, &a); err != nil {
			return fmt.Errorf("failed to unmarshal call: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

// DeleteCall removes an call from the store.
func (s *Store) DeleteCall(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(callsBucket)
		return b.Delete([]byte(id))
	})
}

// AddSentMessage adds a new sent message to the store.
func (s *Store) AddSentMessage(sm *SentMessage) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(sentMessagesBucket)
		id, _ := b.NextSequence()
		sm.ID = fmt.Sprintf("%d", id)

		buf, err := json.Marshal(sm)
		if err != nil {
			return fmt.Errorf("failed to marshal sent message: %w", err)
		}
		return b.Put([]byte(sm.ID), buf)
	})
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

// ListSentMessagesByCallID retrieves all sent messages for a given call ID from the store.
func (s *Store) ListSentMessagesByCallID(callID string) ([]*SentMessage, error) {
	var sentMessages []*SentMessage
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(sentMessagesBucket)
		return b.ForEach(func(k, v []byte) error {
			var sm SentMessage
			if err := json.Unmarshal(v, &sm); err != nil {
				return fmt.Errorf("failed to unmarshal sent message: %w", err)
			}
			if sm.CallID == callID {
				sentMessages = append(sentMessages, &sm)
			}
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list sent messages: %w", err)
	}
	return sentMessages, nil
}

// DeleteSentMessage removes a sent message from the store.
func (s *Store) DeleteSentMessage(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(sentMessagesBucket)
		return b.Delete([]byte(id))
	})
}
