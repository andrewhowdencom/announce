package datastore

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test.db")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	store, err := NewTestStore(tmpfile.Name())
	assert.NoError(t, err)
	defer store.Close()

	// Test AddSentMessage
	sm1 := &SentMessage{
		SourceID:    "1",
		ScheduledAt: time.Now(),
		Timestamp:   "1234567890.123456",
		Status:      StatusSent,
		Type:        "slack",
		Destination: "C1234567890",
	}
	err = store.AddSentMessage("campaign-1", "call-1", sm1)
	assert.NoError(t, err)

	// Test HasBeenSent
	hasBeenSent, err := store.HasBeenSent("campaign-1", "call-1", sm1.Type, sm1.Destination)
	assert.NoError(t, err)
	assert.True(t, hasBeenSent)

	hasBeenSent, err = store.HasBeenSent("campaign-2", "call-2", "slack", "C1234567890")
	assert.NoError(t, err)
	assert.False(t, hasBeenSent)

	// Test ListSentMessages
	sm2 := &SentMessage{
		SourceID:    "2",
		ScheduledAt: time.Now(),
		Timestamp:   "1234567890.123457",
		Status:      StatusSent,
		Type:        "email",
		Destination: "test@example.com",
	}
	err = store.AddSentMessage("campaign-2", "call-2", sm2)
	assert.NoError(t, err)

	sentMessages, err := store.ListSentMessages()
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 2)

	// Test DeleteSentMessage
	err = store.DeleteSentMessage(sm1.ID)
	assert.NoError(t, err)

	hasBeenSent, err = store.HasBeenSent("campaign-1", "call-1", sm1.Type, sm1.Destination)
	assert.NoError(t, err)
	assert.True(t, hasBeenSent)
}
