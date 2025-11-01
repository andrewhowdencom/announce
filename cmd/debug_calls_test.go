package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDebugCallsCmd(t *testing.T) {
	// Create a temporary directory for test files.
	tmpDir, err := ioutil.TempDir("", "ruf-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a valid YAML file with some calls.
	validYAML := `
calls:
  - id: call-1
    destinations:
      - type: slack
        to: ["#general"]
    content: "Hello, world!"
    triggers:
      - scheduled_at: "2025-01-01T12:00:00Z"
  - id: call-2
    destinations:
      - type: slack
        to: ["#random"]
    content: "This is a test."
    triggers:
      - cron: "0 * * * *"
`
	validFile := filepath.Join(tmpDir, "valid.yaml")
	err = ioutil.WriteFile(validFile, []byte(validYAML), 0644)
	assert.NoError(t, err)

	// Create a non-existent file path
	nonExistentFile := filepath.Join(tmpDir, "non-existent.yaml")

	// Set up viper to use a config that points to these two files.
	viper.Set("source.urls", []string{
		"file://" + validFile,
		"file://" + nonExistentFile,
	})

	// Execute the `debug calls` command, capturing stdout and stderr.
	var stdout, stderr bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stderr)
	rootCmd.SetArgs([]string{"debug", "calls"})
	err = rootCmd.Execute()
	assert.NoError(t, err)

	// Assert that stdout contains the correct JSON output for the valid calls.
	expectedJSON := `[
		{
			"id": "call-1",
			"destinations": [
				{
					"type": "slack",
					"to": ["#general"]
				}
			],
			"content": "Hello, world!",
			"triggers": [
				{
					"scheduled_at": "2025-01-01T12:00:00Z"
				}
			],
			"campaign": {
				"id": "valid",
				"name": "` + validFile + `"
			}
		},
		{
			"id": "call-2",
			"destinations": [
				{
					"type": "slack",
					"to": ["#random"]
				}
			],
			"content": "This is a test.",
			"triggers": [
				{
					"cron": "0 * * * *",
					"scheduled_at": "0001-01-01T00:00:00Z"
				}
			],
			"campaign": {
				"id": "valid",
				"name": "` + validFile + `"
			}
		}
	]
	`
	assert.JSONEq(t, expectedJSON, stdout.String())

	// Assert that stderr contains an error message for the non-existent file.
	assert.Contains(t, stderr.String(), "Error sourcing from file://"+nonExistentFile)
}
