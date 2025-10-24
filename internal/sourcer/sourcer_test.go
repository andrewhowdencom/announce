package sourcer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompositeFetcher(t *testing.T) {
	// Test HTTP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", "test-etag")
		fmt.Fprintln(w, "Hello, client")
	}))
	defer server.Close()

	fetcher := NewCompositeFetcher()
	fetcher.AddFetcher("http", NewHTTPFetcher())

	data, state, err := fetcher.Fetch(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, client\n", string(data))
	assert.Equal(t, "test-etag", state)

	// Test File
	tmpfile, err := os.CreateTemp("", "example")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString("Hello, file")
	assert.NoError(t, err)
	tmpfile.Close()

	fetcher.AddFetcher("file", NewFileFetcher())
	fileURL := "file://" + tmpfile.Name()
	data, _, err = fetcher.Fetch(fileURL)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, file", string(data))

	// Test Unsupported Scheme
	_, _, err = fetcher.Fetch("ftp://example.com")
	assert.Error(t, err)
}
