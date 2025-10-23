package slack

import (
	"testing"
	"time"

	"github.com/andrewhowdencom/ruf/internal/datastore"
	"github.com/stretchr/testify/assert"
)

func TestDoReschedule(t *testing.T) {
	// Create a new in-memory datastore for testing
	store, err := datastore.NewMockStore()
	assert.NoError(t, err)

	// Create a new call
	call := &datastore.Call{
		Content:     "Hello, world!",
		ChannelID:   "C1234567890",
		ScheduledAt: time.Now().Add(1 * time.Hour),
		Status:      datastore.StatusPending,
	}
	err = store.AddCall(call)
	assert.NoError(t, err)

	// Reschedule the call
	newScheduledAt := time.Now().Add(2 * time.Hour).UTC()
	err = doReschedule(store, call.ID, newScheduledAt.Format(time.RFC3339))
	assert.NoError(t, err)

	// Get the call from the datastore
	updatedCall, err := store.GetCall(call.ID)
	assert.NoError(t, err)

	// Check that the scheduled time and status have been updated
	assert.Equal(t, newScheduledAt.Truncate(time.Second), updatedCall.ScheduledAt.Truncate(time.Second))
	assert.Equal(t, datastore.StatusPending, updatedCall.Status)
}
