package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)



// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "know installed version of gencli",
	Long:  `This command will help you to know the installed version of gencli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.3.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
