package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

// resetViper clears the global Viper instance state.
func resetViper() {
	viper.Reset()
}

// TestSetDefaultConfig verifies that the default config file is created.
func TestSetDefaultConfig(t *testing.T) {
	resetViper()

	// Create a temporary home directory for the test.
	tempHome := t.TempDir()
	// On Windows, os.UserHomeDir() may return USERPROFILE so set both.
	err := os.Setenv("HOME", tempHome)
	err = os.Setenv("USERPROFILE", tempHome)
	if err != nil {
		t.Fatal(err)
	}
	// Call the function that should create the config file.
	SetDefaultConfig()

	configFilePath := tempHome + "/" + configFileDir + "/" + configFileName + "." + configFileType
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configFilePath)
	}
}

// TestUpdateConfig verifies that updating the config works.
func TestUpdateConfig(t *testing.T) {
	resetViper()

	tempHome := t.TempDir()
	err := os.Setenv("HOME", tempHome)
	err = os.Setenv("USERPROFILE", tempHome)
	if err != nil {
		t.Fatal(err)
	}

	// First, create the default config.
	SetDefaultConfig()
	// Then update a key.
	UpdateConfig("test_key", "test_value")

	// Prepare Viper to read the config file.
	resetViper()
	viper.SetConfigFile(tempHome + "/" + configFileDir + "/" + configFileName + "." + configFileType)
	viper.SetConfigType(configFileType)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	value := viper.GetString("test_key")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}
}

// TestGetConfig verifies that GetConfig returns the correct value.
func TestGetConfig(t *testing.T) {
	resetViper()

	tempHome := t.TempDir()
	err := os.Setenv("HOME", tempHome)
	err = os.Setenv("USERPROFILE", tempHome)
	if err != nil {
		t.Fatal(err)
	}

	SetDefaultConfig()
	UpdateConfig("another_key", "another_value")

	value := GetConfig("another_key")
	if value != "another_value" {
		t.Errorf("Expected 'another_value', got '%s'", value)
	}
}

// TestCheckAPIKey verifies that CheckAPIKey prints a message when the API key is not set.
func TestCheckAPIKey(t *testing.T) {
	resetViper()

	// Ensure GEMINI_API_KEY is not set.
	err := os.Unsetenv("GEMINI_API_KEY")
	if err != nil {
		t.Fatal(err)
	}

	// Capture stdout.
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	CheckAPIKey()

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	os.Stdout = oldStdout

	output := buf.String()
	if !strings.Contains(output, "Please set the GEMINI_API_KEY environment variable") {
		t.Errorf("Expected message about GEMINI_API_KEY, got: %s", output)
	}
}

// TestGetHomeDir verifies that getHomeDir returns a non-empty string.
func TestGetHomeDir(t *testing.T) {
	homeDir := getHomeDir()
	if homeDir == "" {
		t.Fatalf("getHomeDir returned an empty string")
	}
}
