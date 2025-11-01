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

func TestYAMLParser(t *testing.T) {
	// Test with campaign
	yamlWithCampaign := `
campaign:
  id: "test-campaign"
  name: "Test Campaign"
calls:
  - id: "test-call"
    subject: "Test Subject"
    content: "Test Content"
`
	parser := NewYAMLParser()
	source, err := parser.Parse("file:///test.yaml", []byte(yamlWithCampaign))
	assert.NoError(t, err)
	assert.Len(t, source.Calls, 1)
	assert.Equal(t, "test-campaign", source.Calls[0].Campaign.ID)
	assert.Equal(t, "Test Campaign", source.Calls[0].Campaign.Name)

	// Test without campaign
	yamlWithoutCampaign := `
calls:
  - id: "test-call"
    subject: "Test Subject"
    content: "Test Content"
`
	source, err = parser.Parse("file:///test.yaml", []byte(yamlWithoutCampaign))
	assert.NoError(t, err)
	assert.Len(t, source.Calls, 1)
	assert.Equal(t, "test", source.Calls[0].Campaign.ID)
	assert.Equal(t, "/test.yaml", source.Calls[0].Campaign.Name)
}
