package datastore

import (
	"fmt"
	"sync"
	"time"
)

// MockStore is a mock implementation of the Storer interface.
type MockStore struct {
	sentMessages map[string]*SentMessage
	mu           sync.Mutex
}

// NewMockStore creates a new MockStore.
func NewMockStore() *MockStore {
	return &MockStore{
		sentMessages: make(map[string]*SentMessage),
	}
}

// AddSentMessage adds a new sent message to the mock store.
func (s *MockStore) AddSentMessage(sm *SentMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sm.ID = fmt.Sprintf("%s@%s", sm.SourceID, sm.ScheduledAt.Format(time.RFC3339Nano))
	s.sentMessages[sm.ID] = sm
	return nil
}

// HasBeenSent checks if a message with the given sourceID and scheduledAt time has been sent.
func (s *MockStore) HasBeenSent(sourceID string, scheduledAt time.Time) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := fmt.Sprintf("%s@%s", sourceID, scheduledAt.Format(time.RFC3339Nano))
	_, ok := s.sentMessages[id]
	return ok, nil
}

// ListSentMessages retrieves all sent messages from the mock store.
func (s *MockStore) ListSentMessages() ([]*SentMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var sentMessages []*SentMessage
	for _, sm := range s.sentMessages {
		sentMessages = append(sentMessages, sm)
	}
	return sentMessages, nil
}

// DeleteSentMessage removes a sent message from the mock store.
func (s *MockStore) DeleteSentMessage(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sentMessages, id)
	return nil
}

// Close is a no-op for the mock store.
func (s *MockStore) Close() error {
	return nil
}
