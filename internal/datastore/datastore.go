package datastore

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adrg/xdg"
	"go.etcd.io/bbolt"
)

// The name of the bucket where announcements will be stored.
var announcementsBucket = []byte("announcements")

// Status represents the status of an announcement.
type Status string

const (
	// StatusPending means the announcement has been scheduled but not yet sent.
	StatusPending Status = "pending"
	// StatusSent means the announcement has been successfully sent.
	StatusSent Status = "sent"
	// StatusFailed means the announcement failed to send.
	StatusFailed Status = "failed"
	// StatusDeleted means the announcement has been deleted.
	StatusDeleted Status = "deleted"
)

// Announcement represents a message to be sent to a destination.
type Announcement struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	ChannelID   string    `json:"channel_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      Status    `json:"status"`
	Timestamp   string    `json:"timestamp,omitempty"`
}

// Storer is an interface that defines the methods for interacting with the datastore.
type Storer interface {
	AddAnnouncement(a *Announcement) error
	GetAnnouncement(id string) (*Announcement, error)
	ListAnnouncements() ([]*Announcement, error)
	UpdateAnnouncement(a *Announcement) error
	DeleteAnnouncement(id string) error
	Close() error
}

// Store manages the persistence of announcements.
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
		_, err := tx.CreateBucketIfNotExists(announcementsBucket)
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

// AddAnnouncement adds a new announcement to the store.
func (s *Store) AddAnnouncement(a *Announcement) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(announcementsBucket)
		id, _ := b.NextSequence()
		a.ID = fmt.Sprintf("%d", id)

		buf, err := json.Marshal(a)
		if err != nil {
			return fmt.Errorf("failed to marshal announcement: %w", err)
		}
		return b.Put([]byte(a.ID), buf)
	})
}

// ListAnnouncements retrieves all announcements from the store.
func (s *Store) ListAnnouncements() ([]*Announcement, error) {
	var announcements []*Announcement
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(announcementsBucket)
		return b.ForEach(func(k, v []byte) error {
			var a Announcement
			if err := json.Unmarshal(v, &a); err != nil {
				return fmt.Errorf("failed to unmarshal announcement: %w", err)
			}
			announcements = append(announcements, &a)
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list announcements: %w", err)
	}
	return announcements, nil
}

// UpdateAnnouncement updates an existing announcement in the store.
func (s *Store) UpdateAnnouncement(a *Announcement) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(announcementsBucket)
		buf, err := json.Marshal(a)
		if err != nil {
			return fmt.Errorf("failed to marshal announcement: %w", err)
		}
		return b.Put([]byte(a.ID), buf)
	})
}

// GetAnnouncement retrieves a single announcement from the store.
func (s *Store) GetAnnouncement(id string) (*Announcement, error) {
	var a *Announcement
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(announcementsBucket)
		v := b.Get([]byte(id))
		if v == nil {
			return ErrNotFound
		}
		if err := json.Unmarshal(v, &a); err != nil {
			return fmt.Errorf("failed to unmarshal announcement: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

// DeleteAnnouncement removes an announcement from the store.
func (s *Store) DeleteAnnouncement(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(announcementsBucket)
		return b.Delete([]byte(id))
	})
}
