package datastore

import (
	"strconv"
)

// MockStore is a mock implementation of the Store for testing.
type MockStore struct {
	calls map[string]*Call
	sentMessages  map[string]*SentMessage
	nextID        int
}

// NewMockStore creates a new MockStore.
func NewMockStore() (Storer, error) {
	return &MockStore{
		calls: make(map[string]*Call),
		sentMessages:  make(map[string]*SentMessage),
		nextID:        1,
	}, nil
}

// AddCall adds a new call to the mock store.
func (m *MockStore) AddCall(a *Call) error {
	id := strconv.Itoa(m.nextID)
	m.nextID++
	a.ID = id
	m.calls[id] = a
	return nil
}

// GetCall retrieves an call from the mock store.
func (m *MockStore) GetCall(id string) (*Call, error) {
	a, ok := m.calls[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

// ListCalls retrieves all calls from the mock store.
func (m *MockStore) ListCalls() ([]*Call, error) {
	calls := make([]*Call, 0, len(m.calls))
	for _, a := range m.calls {
		calls = append(calls, a)
	}
	return calls, nil
}

// UpdateCall updates an existing call in the mock store.
func (m *MockStore) UpdateCall(a *Call) error {
	m.calls[a.ID] = a
	return nil
}

// DeleteCall removes an call from the mock store.
func (m *MockStore) DeleteCall(id string) error {
	delete(m.calls, id)
	return nil
}

// AddSentMessage adds a new sent message to the mock store.
func (m *MockStore) AddSentMessage(sm *SentMessage) error {
	id := strconv.Itoa(m.nextID)
	m.nextID++
	sm.ID = id
	m.sentMessages[id] = sm
	return nil
}

// ListSentMessages retrieves all sent messages from the mock store.
func (m *MockStore) ListSentMessages() ([]*SentMessage, error) {
	sentMessages := make([]*SentMessage, 0, len(m.sentMessages))
	for _, sm := range m.sentMessages {
		sentMessages = append(sentMessages, sm)
	}
	return sentMessages, nil
}

// ListSentMessagesByCallID retrieves all sent messages for a given call ID from the mock store.
func (m *MockStore) ListSentMessagesByCallID(callID string) ([]*SentMessage, error) {
	sentMessages := make([]*SentMessage, 0)
	for _, sm := range m.sentMessages {
		if sm.CallID == callID {
			sentMessages = append(sentMessages, sm)
		}
	}
	return sentMessages, nil
}

// DeleteSentMessage removes a sent message from the mock store.
func (m *MockStore) DeleteSentMessage(id string) error {
	delete(m.sentMessages, id)
	return nil
}

// Close does nothing for the mock store.
func (m *MockStore) Close() error {
	return nil
}
