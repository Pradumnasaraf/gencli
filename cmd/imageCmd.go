package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var (
	imageFilePath   string
	imageFileFormat string
)

var imageCmd = &cobra.Command{
	Use: "image [your question] --path [image path] --format [image format]",
	Example: "gencli image 'What this image is about?' --path cat.png --format png",
	
	Short: "Know details about an image (Please put your question in quotes)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := imageFunc(args)
		fmt.Println(res)
	},
}

func imageFunc(args []string) string {
	userArgs := strings.Join(args[0:], " ")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	CheckNilError(err)
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	imgData, err := os.ReadFile(imageFilePath)
	CheckNilError(err)

	// Supports multiple image inputs
	prompt := []genai.Part{
		genai.ImageData(imageFileFormat, imgData),
		genai.Text(userArgs),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	CheckNilError(err)

	finalResponse := resp.Candidates[0].Content.Parts[0]
	return fmt.Sprint(finalResponse)

}

func init() {
	imageCmd.Flags().StringVarP(&imageFilePath, "path", "p", "", "Enter the image path")
	imageCmd.Flags().StringVarP(&imageFileFormat, "format", "f", "", "Enter the image format (jpeg, png, etc.)")
	errPathF := imageCmd.MarkFlagRequired("path")
	errFormatF := imageCmd.MarkFlagRequired("format")
	CheckNilError(errPathF)
	CheckNilError(errFormatF)
}
