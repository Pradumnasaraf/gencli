package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var (
	numWords       string
	outputLanguage string
	temperature    float32
	saveOutput     bool
	outputFile     string
)

var searchCmd = &cobra.Command{
	Use:     "search [your question]",
	Example: "gencli search 'What is new in Golang?'",
	Short:   "Ask a question and get a response (Please put your question in quotes)",
	Long:    "Ask a question and get a response in a specified number of words. The default number of words is 150. You can change the number of words by using the --words flag.",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := getApiResponse(args)

		if saveOutput {
			// Create directory if it doesn't exist
			dir := filepath.Dir(outputFile)
			if dir != "." {
				err := os.MkdirAll(dir, 0755)
				CheckNilError(err)
			}

			err := os.WriteFile(outputFile, []byte(res), 0644)
			CheckNilError(err)
			fmt.Printf("Response saved to: %s\n", outputFile)
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
	model.SetTemperature(temperature)
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
	searchCmd.Flags().Float32VarP(&temperature, "temperature", "t", 0.5, "Response creativity (0.0-1.0)")
	searchCmd.Flags().BoolVarP(&saveOutput, "save", "s", false, "Save the output to a file")
	searchCmd.Flags().StringVarP(&outputFile, "output", "o", "output.txt", "Output file name")
}
