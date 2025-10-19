package slack

import (
	"testing"
	"time"

	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/stretchr/testify/assert"
)

func TestDoReschedule(t *testing.T) {
	// Create a new in-memory datastore for testing
	store, err := datastore.NewMockStore()
	assert.NoError(t, err)

	// Create a new announcement
	announcement := &datastore.Announcement{
		Content:     "Hello, world!",
		ChannelID:   "C1234567890",
		ScheduledAt: time.Now().Add(1 * time.Hour),
		Status:      datastore.StatusPending,
	}
	err = store.AddAnnouncement(announcement)
	assert.NoError(t, err)

	// Reschedule the announcement
	newScheduledAt := time.Now().Add(2 * time.Hour).UTC()
	err = doReschedule(store, announcement.ID, newScheduledAt.Format(time.RFC3339))
	assert.NoError(t, err)

	// Get the announcement from the datastore
	updatedAnnouncement, err := store.GetAnnouncement(announcement.ID)
	assert.NoError(t, err)

	// Check that the scheduled time and status have been updated
	assert.Equal(t, newScheduledAt.Truncate(time.Second), updatedAnnouncement.ScheduledAt.Truncate(time.Second))
	assert.Equal(t, datastore.StatusPending, updatedAnnouncement.Status)
}
