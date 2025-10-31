package sourcer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitFetcher(t *testing.T) {
	// This test requires an internet connection to a public repo.
	fetcher := NewGitFetcher()
	data, state, err := fetcher.Fetch("git://github.com/golang/go/tree/master/LICENSE")
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.NotEmpty(t, state)

	// Test with a file in a subdirectory
	data, state, err = fetcher.Fetch("git://github.com/golang/go/tree/master/README.md")
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.NotEmpty(t, state)
}
