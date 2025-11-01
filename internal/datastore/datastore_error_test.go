package datastore

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_ErrorCases(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test.db")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	store, err := NewTestStore(tmpfile.Name())
	assert.NoError(t, err)
	defer store.Close()

	// Test GetSentMessage with a non-existent ID
	_, err = store.GetSentMessage("non-existent-id")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))

	// Test DeleteSentMessage with a non-existent ID
	err = store.DeleteSentMessage("non-existent-id")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}
