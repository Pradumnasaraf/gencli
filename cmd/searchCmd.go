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

var numWords string = "150"

var searchCmd = &cobra.Command{
	Use:   "search [your question]",
	Short: "Ask a question and get a response",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := getApiResponse(args)
		fmt.Println(res)
	},
}

func getApiResponse(args []string) string {
	userArgs := strings.Join(args[0:], " ")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	checkNilError(err)
	defer client.Close()

	// Validate user input is a number
	_, err = strconv.Atoi(numWords)
	if err != nil {
		log.Fatal("Invalid number of words")
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(userArgs+" in "+numWords+" words."))
	checkNilError(err)

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
}

func checkNilError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
