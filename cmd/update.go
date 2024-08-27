package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

}

func update() {
	cmd := exec.Command("go", "install", "github.com/Pradumnasaraf/gencli@latest")
	output, err := cmd.Output()

	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	fmt.Println(string(output))
}
