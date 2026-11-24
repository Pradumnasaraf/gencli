package cmd

import (
	"context"
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

	models := getModelsFunc(context.Background())

	options := make([]string, len(models))
	for i, m := range models {
		options[i] = m.DisplayName
	}

	var selected string
	prompt := &survey.Select{
		Message: "Choose an model:",
		Options: options,
	}

	err := surveyAskOne(prompt, &selected)
	CheckNilError(err)

	// Map selected display name to model ID
	modelID := defaultModels[3].ModelID // Default to Gemini 2.5 Pro
	for _, m := range models {
		if m.DisplayName == selected {
			modelID = m.ModelID
			break
		}
	}

	UpdateConfigFunc("genai_model", modelID)
	fmt.Println("Model updated to:", modelID)
}

func init() {
	rootCmd.AddCommand(changeModelCmd)
}
