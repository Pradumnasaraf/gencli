package test

import (
	"os/exec"
	"strings"
	"testing"
)

// TestRootCmd tests the root command (candy)
func TestImageCmd(t *testing.T) {

	expectedOutput := "A CLI tool to interact with the Gemini API."
	cmd := exec.Command("./gencli image 'What this image is about?' --path cat.png --format png")

	// Capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}

	// Validate the cli output
	got := strings.TrimSpace(string(output)[:43])
	if got != expectedOutput {
		t.Errorf("expected %v, but got: %v", expectedOutput, got)
	}

}