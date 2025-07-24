package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const CliVersion = "v1.10.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Know the installed version of gencli",
	Long:  `This command will help you to know the installed version of gencli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gencli version:", CliVersion, "\nTo check and update the latest version, run 'gencli update'")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
