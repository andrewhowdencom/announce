package sourcer

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitFetcher(t *testing.T) {
	t.Run("public repo", func(t *testing.T) {
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
	})

	t.Run("private repo with auth", func(t *testing.T) {
		viper.Set("git.auth.github.com.username", "testuser")
		viper.Set("git.auth.github.com.token", "testtoken")
		defer viper.Set("git.auth.github.com.username", "")
		defer viper.Set("git.auth.github.com.token", "")

		fetcher := NewGitFetcher()
		// This will fail because the credentials are fake, but it proves that the auth is being used.
		_, _, err := fetcher.Fetch("git://github.com/golang/go/tree/master/LICENSE")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "authentication required")
	})
}
