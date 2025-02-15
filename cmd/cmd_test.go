package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"runtime"
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

// TestUpdateCommand tests the 'update' subcommand which is responsible for updating the CLI.
// It uses a mocked execCommand to simulate both successful and failing update scenarios.
func TestUpdateCommand(t *testing.T) {
	// Backup the original execCommand function.
	originalExecCommand := execCommand
	// Restore execCommand after the test completes.
	defer func() { execCommand = originalExecCommand }()

	t.Run("successful_update", func(t *testing.T) {
		// Override execCommand to simulate a successful update.
		execCommand = func(name string, args ...string) *exec.Cmd {
			if runtime.GOOS == "windows" {
				// On Windows, use "cmd /C echo" to simulate a successful response.
				cmdArgs := append([]string{"/C", "echo", "successful update"}, args...)
				cmd := exec.Command("cmd", cmdArgs...)
				cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
				return cmd
			}
			// On Unix-like systems, directly use echo.
			cmd := exec.Command("echo", append([]string{"successful update"}, args...)...)
			cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
			return cmd
		}

		// Execute the update command.
		output, err := executeCommand(t, rootCmd, "update")
		// Verify that the command executed without error.
		assert.NoError(t, err)
		// Confirm that the output contains a success message.
		assert.Contains(t, output, "CLI updated successfully")
	})

	t.Run("failed_update", func(t *testing.T) {
		// Override execCommand to simulate a failing update.
		execCommand = func(name string, args ...string) *exec.Cmd {
			if runtime.GOOS == "windows" {
				// Simulate failure on Windows by exiting with a non-zero status.
				return exec.Command("cmd", "/C", "exit", "1")
			}
			// On Unix-like systems, use the false command which always fails.
			return exec.Command("false")
		}

		// Execute the update command and expect an error.
		_, err := executeCommand(t, rootCmd, "update")
		assert.Error(t, err)
	})
}

// TestSearchCommand tests the 'search' subcommand using table-driven tests.
// It covers different scenarios like basic search, empty query, different languages, and API errors.
func TestSearchCommand(t *testing.T) {
	// Backup the original getApiResponseFunc.
	originalFunc := getApiResponseFunc
	// Restore the original function after tests.
	defer func() { getApiResponseFunc = originalFunc }()

	// Define test cases for the search command.
	testCases := []struct {
		name           string   // Name of the test case.
		args           []string // Command line arguments to pass.
		mockResponse   string   // The response to simulate from the API.
		expectedOutput string   // Expected output to be verified.
	}{
		{
			name:           "basic_search",
			args:           []string{"search", "test query", "--words", "100"},
			mockResponse:   "test response",
			expectedOutput: "test response",
		},
		{
			name: "empty_query",
			// Cobra requires at least one argument; simulate an empty query with an empty string.
			args:           []string{"search", ""},
			mockResponse:   "query cannot be empty",
			expectedOutput: "query cannot be empty",
		},
		{
			name:           "different_language",
			args:           []string{"search", "query", "--language", "french"},
			mockResponse:   "réponse en français",
			expectedOutput: "réponse en français",
		},
		{
			name:           "error_response",
			args:           []string{"search", "error"},
			mockResponse:   "API_ERROR",
			expectedOutput: "API_ERROR",
		},
	}

	// Loop through each test case.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Override getApiResponseFunc to return the mock response for the test.
			getApiResponseFunc = func(args []string) string {
				// If the query is empty, return an appropriate error message.
				if len(args) > 0 && args[0] == "" {
					return "query cannot be empty"
				}
				return tc.mockResponse
			}

			// Execute the search command with provided arguments.
			output, err := executeCommand(t, rootCmd, tc.args...)
			// The search command should not return an error; it prints the result instead.
			assert.NoError(t, err)
			// Check that the output contains the expected API response.
			assert.Contains(t, output, tc.expectedOutput)
		})
	}
}
