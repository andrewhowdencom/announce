package datastore

// MockStore is a mock implementation of the Store for testing.
type MockStore struct {
	AddAnnouncementFunc    func(*Announcement) error
	ListAnnouncementsFunc  func() ([]*Announcement, error)
	UpdateAnnouncementFunc func(*Announcement) error
	DeleteAnnouncementFunc func(string) error
	CloseFunc              func() error
}

// AddAnnouncement calls the AddAnnouncementFunc.
func (m *MockStore) AddAnnouncement(a *Announcement) error {
	return m.AddAnnouncementFunc(a)
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

// Close calls the CloseFunc.
func (m *MockStore) Close() error {
	return m.CloseFunc()
}
