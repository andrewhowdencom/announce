/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package slack

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPostCmd(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request body
		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"text": "Hello, world!"}`, string(body))

		// Send a response
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Set the webhook URL
	viper.Set("slack.webhook_url", server.URL)

	// Redirect STDIN
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.Write([]byte("Hello, world!"))
	w.Close()

	// Redirect STDOUT
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	r, w, _ = os.Pipe()
	os.Stdout = w

	// Run the command
	err := PostCmd.RunE(PostCmd, []string{})
	assert.NoError(t, err)

	// Check the output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Equal(t, "Message sent to Slack successfully\n", buf.String())
}