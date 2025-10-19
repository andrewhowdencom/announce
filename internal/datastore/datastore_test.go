package datastore

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	// Create a new store for testing
	store, err := NewStore()
	assert.NoError(t, err)
	defer store.Close()
	defer os.RemoveAll(".config/announce")

	// Test AddAnnouncement
	a := &Announcement{
		Content:     "Hello, world!",
		ChannelID:   "C1234567890",
		ScheduledAt: time.Now(),
		Status:      StatusPending,
	}
	err = store.AddAnnouncement(a)
	assert.NoError(t, err)
	assert.NotEmpty(t, a.ID)

	// Test ListAnnouncements
	announcements, err := store.ListAnnouncements()
	assert.NoError(t, err)
	assert.Len(t, announcements, 1)
	assert.Equal(t, a.Content, announcements[0].Content)

	// Test UpdateAnnouncement to StatusSent
	announcements[0].Status = StatusSent
	err = store.UpdateAnnouncement(announcements[0])
	assert.NoError(t, err)

	announcements, err = store.ListAnnouncements()
	assert.NoError(t, err)
	assert.Equal(t, StatusSent, announcements[0].Status)

	// Test UpdateAnnouncement to StatusProcessed
	announcements[0].Status = StatusProcessed
	err = store.UpdateAnnouncement(announcements[0])
	assert.NoError(t, err)

	processedAnnouncement, err := store.GetAnnouncement(announcements[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, StatusProcessed, processedAnnouncement.Status)

	// Test AddSentMessage
	sm := &SentMessage{
		AnnouncementID: a.ID,
		Timestamp:      "12345",
		Status:         StatusSent,
	}
	err = store.AddSentMessage(sm)
	assert.NoError(t, err)
	assert.NotEmpty(t, sm.ID)

	// Test ListSentMessages
	sentMessages, err := store.ListSentMessages()
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 1)
	assert.Equal(t, sm.AnnouncementID, sentMessages[0].AnnouncementID)

	// Test ListSentMessagesByAnnouncementID
	sentMessages, err = store.ListSentMessagesByAnnouncementID(a.ID)
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 1)
	assert.Equal(t, sm.AnnouncementID, sentMessages[0].AnnouncementID)

	// Test DeleteSentMessage
	err = store.DeleteSentMessage(sm.ID)
	assert.NoError(t, err)

	sentMessages, err = store.ListSentMessages()
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 0)

	// Test DeleteAnnouncement
	err = store.DeleteAnnouncement(announcements[0].ID)
	assert.NoError(t, err)

	announcements, err = store.ListAnnouncements()
	assert.NoError(t, err)
	assert.Len(t, announcements, 0)
}
