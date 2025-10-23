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
		fmt.Fprintln(w, "Hello, client")
	}))
	defer server.Close()

	fetcher := NewCompositeFetcher()
	fetcher.AddFetcher("http", NewHTTPFetcher())

	data, err := fetcher.Fetch(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, client\n", string(data))

	// Test File
	tmpfile, err := os.CreateTemp("", "example")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString("Hello, file")
	assert.NoError(t, err)
	tmpfile.Close()

	fetcher.AddFetcher("file", NewFileFetcher())
	fileURL := "file://" + tmpfile.Name()
	data, err = fetcher.Fetch(fileURL)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, file", string(data))

	// Test Unsupported Scheme
	_, err = fetcher.Fetch("ftp://example.com")
	assert.Error(t, err)
}
