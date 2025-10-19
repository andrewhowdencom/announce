package datastore

import (
	"strconv"
)

// MockStore is a mock implementation of the Store for testing.
type MockStore struct {
	announcements map[string]*Announcement
	sentMessages  map[string]*SentMessage
	nextID        int
}

// NewMockStore creates a new MockStore.
func NewMockStore() (Storer, error) {
	return &MockStore{
		announcements: make(map[string]*Announcement),
		sentMessages:  make(map[string]*SentMessage),
		nextID:        1,
	}, nil
}

// AddAnnouncement adds a new announcement to the mock store.
func (m *MockStore) AddAnnouncement(a *Announcement) error {
	id := strconv.Itoa(m.nextID)
	m.nextID++
	a.ID = id
	m.announcements[id] = a
	return nil
}

// GetAnnouncement retrieves an announcement from the mock store.
func (m *MockStore) GetAnnouncement(id string) (*Announcement, error) {
	a, ok := m.announcements[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

// ListAnnouncements retrieves all announcements from the mock store.
func (m *MockStore) ListAnnouncements() ([]*Announcement, error) {
	announcements := make([]*Announcement, 0, len(m.announcements))
	for _, a := range m.announcements {
		announcements = append(announcements, a)
	}
	return announcements, nil
}

// UpdateAnnouncement updates an existing announcement in the mock store.
func (m *MockStore) UpdateAnnouncement(a *Announcement) error {
	m.announcements[a.ID] = a
	return nil
}

// DeleteAnnouncement removes an announcement from the mock store.
func (m *MockStore) DeleteAnnouncement(id string) error {
	delete(m.announcements, id)
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

// ListSentMessagesByAnnouncementID retrieves all sent messages for a given announcement ID from the mock store.
func (m *MockStore) ListSentMessagesByAnnouncementID(announcementID string) ([]*SentMessage, error) {
	sentMessages := make([]*SentMessage, 0)
	for _, sm := range m.sentMessages {
		if sm.AnnouncementID == announcementID {
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
