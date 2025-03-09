package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	// execCommand is a variable for creating commands. It can be overridden in tests.
	execCommand = exec.Command

	//nolint:unused // exitFunc wraps os.Exit so that it can be overridden in tests if needed.
	exitFunc = func(code int) {
		os.Exit(code)
	}
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update gencli to the latest version",
	Long:  `This command will help you to update gencli to the latest version.`,
	// Use RunE so that we can return errors.
	RunE: func(cmd *cobra.Command, args []string) error {
		return update()
	},
}

func update() error {
	// Run the "go install" command to update the CLI.
	cmd := execCommand("go", "install", "github.com/Pradumnasaraf/gencli@latest")
	_, err := cmd.Output()
	if err != nil {
		// Instead of calling CheckNilError (which might exit), return the error.
		return fmt.Errorf("failed to update CLI: %w", err)
	}

	fmt.Printf("CLI updated successfully to the latest version (If any).\n")
	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
