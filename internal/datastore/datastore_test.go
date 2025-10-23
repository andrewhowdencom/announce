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

	// Test AddCall
	a := &Call{
		Content:     "Hello, world!",
		ChannelID:   "C1234567890",
		ScheduledAt: time.Now(),
		Status:      StatusPending,
	}
	err = store.AddCall(a)
	assert.NoError(t, err)
	assert.NotEmpty(t, a.ID)

	// Test ListCalls
	calls, err := store.ListCalls()
	assert.NoError(t, err)
	assert.Len(t, calls, 1)
	assert.Equal(t, a.Content, calls[0].Content)

	// Test UpdateCall to StatusSent
	calls[0].Status = StatusSent
	err = store.UpdateCall(calls[0])
	assert.NoError(t, err)

	calls, err = store.ListCalls()
	assert.NoError(t, err)
	assert.Equal(t, StatusSent, calls[0].Status)

	// Test UpdateCall to StatusProcessed
	calls[0].Status = StatusProcessed
	err = store.UpdateCall(calls[0])
	assert.NoError(t, err)

	processedCall, err := store.GetCall(calls[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, StatusProcessed, processedCall.Status)

	// Test AddSentMessage
	sm := &SentMessage{
		CallID: a.ID,
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
	assert.Equal(t, sm.CallID, sentMessages[0].CallID)

	// Test ListSentMessagesByCallID
	sentMessages, err = store.ListSentMessagesByCallID(a.ID)
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 1)
	assert.Equal(t, sm.CallID, sentMessages[0].CallID)

	// Test DeleteSentMessage
	err = store.DeleteSentMessage(sm.ID)
	assert.NoError(t, err)

	sentMessages, err = store.ListSentMessages()
	assert.NoError(t, err)
	assert.Len(t, sentMessages, 0)

	// Test DeleteCall
	err = store.DeleteCall(calls[0].ID)
	assert.NoError(t, err)

	calls, err = store.ListCalls()
	assert.NoError(t, err)
	assert.Len(t, calls, 0)
}
