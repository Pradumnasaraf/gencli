package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/AlecAivazis/survey/v2"
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
		assert.Contains(t, output, CliVersion)
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

// TestImageCommand tests the 'image' subcommand which analyzes images.
// It covers scenarios like valid image path, invalid path, unsupported image format, and API error.
func TestImageCommand(t *testing.T) {
	// Backup the original getApiResponseImageFunc.
	originalFunc := getApiResponseImageFunc
	// Restore it after the tests.
	defer func() { getApiResponseImageFunc = originalFunc }()

	// Define file paths used in the tests.
	validImagePath := filepath.Join("..", "assets", "test.jpg")
	invalidImagePath := filepath.Join("..", "assets", "missing.jpg")

	// Define test cases.
	testCases := []struct {
		name          string   // Name of the test case.
		args          []string // Command line arguments to pass.
		mockResponse  string   // Response to simulate (for non-error cases).
		expectError   bool     // Whether we expect an error to be shown in the output.
		errorContains string   // Substring that should be present in the error message.
	}{
		{
			name:         "valid_image",
			args:         []string{"image", "test query", "--path", validImagePath, "--format", "jpg", "--language", "english"},
			mockResponse: "image analysis",
		},
		{
			name: "invalid_path",
			// Simulate file-read error using an invalid file path.
			args:          []string{"image", "query", "--path", invalidImagePath, "--format", "jpg", "--language", "english"},
			expectError:   true,
			errorContains: "no such file",
		},
		{
			name: "unsupported_format",
			// Provide an unsupported image format to simulate error.
			args:          []string{"image", "query", "--path", validImagePath, "--format", "bmp", "--language", "english"},
			expectError:   true,
			errorContains: "unsupported format",
		},
		{
			name: "api_error",
			// Simulate an API error by returning "API_ERROR".
			args:          []string{"image", "query", "--path", validImagePath, "--format", "jpg", "--language", "english"},
			mockResponse:  "API_ERROR",
			expectError:   true,
			errorContains: "API_ERROR",
		},
	}

	// Iterate over each test case.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Override getApiResponseImageFunc to simulate different responses based on flags.
			getApiResponseImageFunc = func(args []string) string {
				// Check if the required image path flag was provided.
				if imageFilePath == "" {
					return "Error: required flag \"path\" not set"
				}
				// Simulate error if an invalid file path is provided.
				if imageFilePath == invalidImagePath {
					return "Error: no such file"
				}
				// Simulate error for unsupported image formats.
				if imageFileFormat == "bmp" {
					return "Error: unsupported format"
				}
				// If mockResponse is "API_ERROR", simulate an API error.
				if tc.mockResponse == "API_ERROR" {
					return "API_ERROR"
				}
				// Otherwise, return the provided mock response.
				return tc.mockResponse
			}

			// Execute the image command.
			output, err := executeCommand(t, rootCmd, tc.args...)
			// Even if there is an error in processing, the command itself doesn't return an error (uses Run, not RunE).
			assert.NoError(t, err)
			if tc.expectError {
				// Verify that the error message is included in the output.
				assert.Contains(t, output, tc.errorContains)
			} else {
				// Verify that the successful response is printed.
				assert.Contains(t, output, tc.mockResponse)
			}
		})
	}
}

// TestModelCommand tests the 'model' subcommand which allows the user to change the current model.
// It simulates user selection using the survey package and verifies that the configuration is updated.
func TestModelCommand(t *testing.T) {
	// Backup the original functions to allow restoration later.
	originalSurveyAskOne := surveyAskOne
	originalGetConfigFunc := GetConfigFunc
	originalUpdateConfigFunc := UpdateConfigFunc

	// Restore the original functions after tests.
	defer func() {
		surveyAskOne = originalSurveyAskOne
		GetConfigFunc = originalGetConfigFunc
		UpdateConfigFunc = originalUpdateConfigFunc
	}()

	// Setup an in-memory configuration map for testing purposes.
	testConfig := make(map[string]string)
	GetConfigFunc = func(key string) string {
		return testConfig[key]
	}
	UpdateConfigFunc = func(key, value string) {
		testConfig[key] = value
	}

	// Define test cases for different model selections.
	testCases := []struct {
		name            string // Test case name.
		mockSelection   string // Simulated user input from the survey.
		expectedModel   string // Expected model configuration after the command.
		expectedMessage string // Expected output message indicating model update.
	}{
		{
			name:            "select_gemini_2.0_flash",
			mockSelection:   "Gemini 2.0 Flash",
			expectedModel:   "gemini-2.0-flash",
			expectedMessage: "Model updated to: gemini-2.0-flash",
		},
		{
			name:            "select_gemini_2.0_flash_lite",
			mockSelection:   "Gemini 2.0 Flash-Lite Preview",
			expectedModel:   "gemini-2.0-flash-lite",
			expectedMessage: "Model updated to: gemini-2.0-flash-lite",
		},
		{
			name:            "invalid_selection_default",
			mockSelection:   "Invalid Model",
			expectedModel:   "gemini-2.0-flash", // Default model in case of an invalid selection.
			expectedMessage: "Model updated to: gemini-2.0-flash",
		},
	}

	// Iterate over each model selection test case.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Override surveyAskOne to simulate the user's selection.
			surveyAskOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
				resp := response.(*string)
				*resp = tc.mockSelection
				return nil
			}

			// Set an initial model configuration.
			testConfig["genai_model"] = "gemini-1.5-flash"

			// Execute the model command.
			output, err := executeCommand(t, rootCmd, "model")
			require.NoError(t, err)

			// Verify that the output contains the current model and the update message.
			assert.Contains(t, output, "Current model: gemini-1.5-flash")
			assert.Contains(t, output, tc.expectedMessage)
			// Confirm that the configuration has been updated to the expected model.
			currentModel := GetConfigFunc("genai_model")
			assert.Equal(t, tc.expectedModel, currentModel)
		})
	}

	// Test to ensure that the 'model' command is registered with the root command.
	t.Run("command_registration", func(t *testing.T) {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == "model" {
				found = true
				break
			}
		}
		assert.True(t, found, "model command should be registered")
	})
}

// TestErrorHandling tests how the CLI handles invalid commands.
// It verifies that an appropriate error message is shown when an unknown command is used.
func TestErrorHandling(t *testing.T) {
	t.Run("invalid_command", func(t *testing.T) {
		// Attempt to execute a command that doesn't exist.
		_, err := executeCommand(t, rootCmd, "invalid-command")
		// Verify that an error is returned.
		assert.Error(t, err)
		// Check that the error message mentions "unknown command".
		assert.Contains(t, err.Error(), "unknown command")
	})
}
