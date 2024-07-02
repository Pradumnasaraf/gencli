package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Pradumnasaraf/gencli/config"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var rootCmd = &cobra.Command{
	Use:   "gencli",
	Short: "Ask me anything :)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := getApiRespone(args)
		fmt.Println(res)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getApiRespone(args []string) string {

	// Load the environment variables
	config.Config()

	userArgs := strings.Join(args[1:], " ")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))

	if err != nil {
		if strings.Contains(err.Error(), "GEMINI_API_KEY") {
			log.Fatal("Please set the GEMINI_API_KEY environment variable. Check the README for more information.")
		}
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(userArgs+"in 100-120 words."))
	if err != nil {
		log.Fatal(err)
	}

	finalResponse := resp.Candidates[0].Content.Parts[0]

	return fmt.Sprint(finalResponse)
}
