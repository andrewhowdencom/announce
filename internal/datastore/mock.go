package datastore

import (
	"fmt"
	"strings"
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
	sm.ID = s.generateID(sm.SourceID, sm.ScheduledAt, sm.Type, sm.Destination)
	s.sentMessages[sm.ID] = sm
	return nil
}

// HasBeenSent checks if a message with the given sourceID and scheduledAt time has been sent.
func (s *MockStore) HasBeenSent(sourceID string, scheduledAt time.Time, destType, destination string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.generateID(sourceID, scheduledAt, destType, destination)
	sm, ok := s.sentMessages[id]
	return ok && sm.Status != StatusDeleted, nil
}

func (s *MockStore) generateID(sourceID string, scheduledAt time.Time, destType, destination string) string {
	parts := []string{
		sourceID,
		scheduledAt.Format(time.RFC3339Nano),
		destType,
		destination,
	}
	return strings.Join(parts, "@")
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

// GetSentMessage retrieves a single sent message from the mock store.
func (s *MockStore) GetSentMessage(id string) (*SentMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sm, ok := s.sentMessages[id]
	if !ok {
		return nil, fmt.Errorf("message with id '%s' not found", id)
	}
	return sm, nil
}

// DeleteSentMessage removes a sent message from the mock store.
func (s *MockStore) DeleteSentMessage(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sm, ok := s.sentMessages[id]
	if !ok {
		return fmt.Errorf("message with id '%s' not found", id)
	}
	sm.Status = StatusDeleted
	return nil
}

// Close is a no-op for the mock store.
func (s *MockStore) Close() error {
	return nil
}
