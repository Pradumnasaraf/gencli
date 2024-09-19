package test

import (
	"os/exec"
	"strings"
	"testing"
)

// TestRootCmd tests the root command (gencli)
func TestRootCmd(t *testing.T) {

	expectedOutput := "A CLI tool to interact with the Gemini API."
	cmd := exec.Command("./gencli")

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

// TestRootCmdHelpFlag tests the root command (gencli) with help flag
func TestRootCmdHelpFlag(t *testing.T) {

	expectedOutput := "Do all your tedious tasks with a single command"
	cmd := exec.Command("./gencli", "--help")

	// Capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}

	// Validate the cli output
	got := strings.TrimSpace(string(output)[:47])
	if got != expectedOutput {
		t.Errorf("expected %v, but got: %v", expectedOutput, got)
	}

}
