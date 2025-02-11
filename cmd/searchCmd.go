package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var (
	numWords       string  = "150"
	outputLanguage string  = "english"
	temperature    float32 = 0.7
	saveOutput     string
)

var searchCmd = &cobra.Command{
	Use:     "search [your question]",
	Example: "gencli search 'What is new in Golang?'",
	Short:   "Ask a question and get a response (Please put your question in quotes)",
	Long:    "Ask a question and get a response in a specified number of words. The default number of words is 150. You can change the number of words by using the --words flag.",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := getApiResponse(args)
		filename, _ := cmd.Flags().GetString("save-output")
		if filename == "" {
			filename = "output.txt"
		}
		if cmd.Flags().Changed("save-output") {
			err := os.WriteFile(filename, []byte(res), 0644)
			CheckNilError(err)
			fmt.Printf("Response saved to: %s\n", filename)
		} else {
			fmt.Println(res)
		}
	},
}

func getApiResponse(args []string) string {
	userArgs := strings.Join(args[0:], " ")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	CheckNilError(err)
	defer client.Close()

	// Validate user input is a number
	_, err = strconv.Atoi(numWords)
	if err != nil {
		log.Fatal("Invalid number of words")
	}

	currentGenaiModel := GetConfig("genai_model")
	model := client.GenerativeModel(currentGenaiModel)
	resp, err := model.GenerateContent(ctx, genai.Text(userArgs+" in "+numWords+" words"+" in "+outputLanguage+" language"))
	CheckNilError(err)

	finalResponse := resp.Candidates[0].Content.Parts[0]

	return formatAsPlainText(fmt.Sprint(finalResponse))
}

func formatAsPlainText(input string) string {
	// Remove Markdown headings
	re := regexp.MustCompile(`^#{1,6}\s+(.*)`)
	input = re.ReplaceAllString(input, "$1")

	// Remove horizontal rules
	re = regexp.MustCompile(`\n---\n`)
	input = re.ReplaceAllString(input, "\n")

	// Remove italic formatting
	re = regexp.MustCompile(`_([^_]+)_`)
	input = re.ReplaceAllString(input, "$1")

	// Remove bold formatting
	re = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	input = re.ReplaceAllString(input, "$1")

	// Remove bullet points
	re = regexp.MustCompile(`\n\* (.*)`)
	input = re.ReplaceAllString(input, "\n- $1")

	return input
}

func init() {
	searchCmd.Flags().StringVarP(&numWords, "words", "w", "150", "Number of words in the response")
	searchCmd.Flags().StringVarP(&outputLanguage, "language", "l", "english", "Output language")
	searchCmd.Flags().Float32VarP(&temperature, "temperature", "t", 0.7, "Response creativity (0.0-1.0)")

	// Corrected flag handling
	searchCmd.Flags().StringVarP(&saveOutput, "save-output", "s", "", "Save response to a file (default: output.txt if used without a value)")
	searchCmd.Flags().Lookup("save-output").NoOptDefVal = "output.txt" // If no argument, default to output.txt
}
