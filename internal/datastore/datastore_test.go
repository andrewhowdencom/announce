package datastore

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
)

func TestStore(t *testing.T) {
	// Create a temporary database for testing
	f, err := ioutil.TempFile("", "test.db")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	db, err := bbolt.Open(f.Name(), 0600, nil)
	assert.NoError(t, err)

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(announcementsBucket)
		return err
	})
	assert.NoError(t, err)

	store := &Store{db: db}
	defer store.Close()

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

	// Test UpdateAnnouncement
	announcements[0].Status = StatusSent
	err = store.UpdateAnnouncement(announcements[0])
	assert.NoError(t, err)

	announcements, err = store.ListAnnouncements()
	assert.NoError(t, err)
	assert.Equal(t, StatusSent, announcements[0].Status)

	// Test DeleteAnnouncement
	err = store.DeleteAnnouncement(announcements[0].ID)
	assert.NoError(t, err)

	announcements, err = store.ListAnnouncements()
	assert.NoError(t, err)
	assert.Len(t, announcements, 0)
}
