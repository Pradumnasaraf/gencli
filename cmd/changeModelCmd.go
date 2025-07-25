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

	currentGenaiModel := GetConfigFunc("genai_model")
	fmt.Println("Current model:", currentGenaiModel)

	options := []string{"Gemini 2.5 Pro", "Gemini 2.5 Flash", "Gemini 2.5 Flash-Lite", "Gemini 2.0 Flash", "Gemini 2.0 Flash-Lite"}

	var selected string
	prompt := &survey.Select{
		Message: "Choose an model:",
		Options: options,
	}

	err := surveyAskOne(prompt, &selected)
	CheckNilError(err)

	switch selected {
	case "Gemini 2.5 Pro":
		selected = "gemini-2.5-pro"
	case "Gemini 2.5 Flash":
		selected = "gemini-2.5-flash"
	case "Gemini 2.5 Flash-Lite":
		selected = "gemini-2.5-flash-lite"
	case "Gemini 2.0 Flash":
		selected = "gemini-2.0-flash"
	case "Gemini 2.0 Flash-Lite":
		selected = "gemini-2.0-flash-lite"
	default:
		selected = "gemini-2.5-pro"
	}
	UpdateConfigFunc("genai_model", selected)

	fmt.Println("Model updated to:", selected)
}

func init() {
	rootCmd.AddCommand(changeModelCmd)
}
