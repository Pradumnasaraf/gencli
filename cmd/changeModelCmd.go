package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var changeModelCmd = &cobra.Command{
	Use:   "model",
	Short: "To select a different GenAI model",
	Long:  `This command will help you to select a different GenAI model.`,
	Run: func(cmd *cobra.Command, args []string) {
		setModelConfig()
	},
}

func setModelConfig() {

	currentGenaiModel := GetConfig("genai_model")
	fmt.Println("Current model:", currentGenaiModel)

	options := []string{"Gemini 2.0 Flash", "Gemini 2.0 Flash-Lite Preview", "Gemini 1.5 Flash", "Gemini 1.5 Flash-8B", "Gemini 1.5 Pro"}

	var selected string
	prompt := &survey.Select{
		Message: "Choose an model:",
		Options: options,
	}

	err := survey.AskOne(prompt, &selected)
	if err != nil {
		fmt.Println("Error:", err)
	}

	switch selected {
	case "Gemini 2.0 Flash":
		selected = "gemini-2.0-flash"
	case "Gemini 2.0 Flash-Lite Preview":
		selected = "gemini-2.0-flash-lite-preview"
	case "Gemini 1.5 Flash":
		selected = "gemini-1.5-flash"
	case "Gemini 1.5 Flash-8B":
		selected = "gemini-1.5-flash-8b"
	case "Gemini 1.5 Pro":
		selected = "gemini-1.5-pro"
	default:
		selected = "gemini-1.5-flash"
	}
	UpdateConfig("genai_model", selected)

	fmt.Println("Model updated to:", selected)
}

func init() {
	rootCmd.AddCommand(changeModelCmd)
}
