package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMain sets up the test environment before running any tests.
// It initializes environment variables that are needed by the CLI.
func TestMain(m *testing.M) {
	// Create a dummy testing.T instance to set environment variables.
	t := &testing.T{}
	// Set a dummy API key so that API calls in tests don't fail.
	t.Setenv("GEMINI_API_KEY", "test-key")

	// Run all tests.
	os.Exit(m.Run())
}

// executeCommand is a helper function that runs a cobra command with the given arguments,
// captures the output printed to stdout, and returns it along with any execution error.
// It uses testify assertions to ensure that the pipe setup and IO copying are successful.
func executeCommand(t *testing.T, root *cobra.Command, args ...string) (string, error) {
	t.Helper()

	// Save the current stdout.
	oldStdout := os.Stdout
	// Create a pipe to capture output.
	r, w, err := os.Pipe()
	require.NoError(t, err)

	// Redirect stdout to the write end of the pipe.
	os.Stdout = w
	// Ensure stdout is restored after command execution.
	defer func() { os.Stdout = oldStdout }()

	// Set the command arguments.
	root.SetArgs(args)
	// Execute the command.
	_, cmdErr := root.ExecuteC()

	// Close the write end of the pipe to finish capturing output.
	err = w.Close()
	if err != nil {
		return "", err
	}
	// Copy the output from the read end of the pipe into a buffer.
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	// Return the captured output and any error from command execution.
	return buf.String(), cmdErr
}

// TestVersionCommand tests the 'version' subcommand.
// It verifies that the command prints the expected version string.
func TestVersionCommand(t *testing.T) {
	t.Run("version_output", func(t *testing.T) {
		// Execute the version command and capture its output.
		output, err := executeCommand(t, rootCmd, "version")
		// Check that no error occurred during execution.
		assert.NoError(t, err)
		// Verify that the output contains the expected version (e.g., "v1.6.3").
		assert.Contains(t, output, "v1.6.3")
	})
}
