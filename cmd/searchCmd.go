package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
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
		res := getApiRespone(args)
		fmt.Println(res)
	},
}

func getApiRespone(args []string) string {

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

	return fmt.Sprint(finalResponse)
}

func init() {
	searchCmd.Flags().StringVarP(&numWords, "words", "w", "150", "Number of words in the response")
}

func checkNilError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
