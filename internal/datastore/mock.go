package datastore

// MockStore is a mock implementation of the Store for testing.
type MockStore struct {
	AddAnnouncementFunc                func(*Announcement) error
	GetAnnouncementFunc                func(string) (*Announcement, error)
	ListAnnouncementsFunc              func() ([]*Announcement, error)
	UpdateAnnouncementFunc             func(*Announcement) error
	DeleteAnnouncementFunc             func(string) error
	AddSentMessageFunc                 func(*SentMessage) error
	ListSentMessagesFunc               func() ([]*SentMessage, error)
	ListSentMessagesByAnnouncementIDFunc func(string) ([]*SentMessage, error)
	DeleteSentMessageFunc              func(string) error
	CloseFunc                          func() error
}

// AddAnnouncement calls the AddAnnouncementFunc.
func (m *MockStore) AddAnnouncement(a *Announcement) error {
	return m.AddAnnouncementFunc(a)
}

// GetAnnouncement calls the GetAnnouncementFunc.
func (m *MockStore) GetAnnouncement(id string) (*Announcement, error) {
	return m.GetAnnouncementFunc(id)
}

// ListAnnouncements calls the ListAnnouncementsFunc.
func (m *MockStore) ListAnnouncements() ([]*Announcement, error) {
	return m.ListAnnouncementsFunc()
}

// UpdateAnnouncement calls the UpdateAnnouncementFunc.
func (m *MockStore) UpdateAnnouncement(a *Announcement) error {
	return m.UpdateAnnouncementFunc(a)
}

// DeleteAnnouncement calls the DeleteAnnouncementFunc.
func (m *MockStore) DeleteAnnouncement(id string) error {
	return m.DeleteAnnouncementFunc(id)
}

// AddSentMessage calls the AddSentMessageFunc.
func (m *MockStore) AddSentMessage(sm *SentMessage) error {
	return m.AddSentMessageFunc(sm)
}

// ListSentMessages calls the ListSentMessagesFunc.
func (m *MockStore) ListSentMessages() ([]*SentMessage, error) {
	return m.ListSentMessagesFunc()
}

// ListSentMessagesByAnnouncementID calls the ListSentMessagesByAnnouncementIDFunc.
func (m *MockStore) ListSentMessagesByAnnouncementID(announcementID string) ([]*SentMessage, error) {
	return m.ListSentMessagesByAnnouncementIDFunc(announcementID)
}

// DeleteSentMessage calls the DeleteSentMessageFunc.
func (m *MockStore) DeleteSentMessage(id string) error {
	return m.DeleteSentMessageFunc(id)
}

// Close calls the CloseFunc.
func (m *MockStore) Close() error {
	return m.CloseFunc()
}
