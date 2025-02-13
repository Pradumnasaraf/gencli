package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var (
	imageFilePath      string
	imageFileFormat    string
	respOutputLanguage string
	saveResponse       bool
	saveResponseFile   string
	modelTemp          float32
)

var imageCmd = &cobra.Command{
	Use:     "image [your question] --path [image path] --format [image format] --language [output language] --temperature [creativity] --save --output [output file]",
	Example: "gencli image 'What this image is about?' --path cat.png --format png",
	Short:   "Know details about an image (Please put your question in quotes)",
	Long:    "Ask a question about an image and get a response. You need to provide the path of the image and the format of the image. The supported formats are jpg, jpeg, png, and gif.",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := getApiResponseImage(args)
		if saveResponse {
			// Create directory if it doesn't exist
			dir := filepath.Dir(saveResponseFile)
			if dir != "." {
				err := os.MkdirAll(dir, 0755)
				CheckNilError(err)
			}

			err := os.WriteFile(saveResponseFile, []byte(res), 0644)
			CheckNilError(err)
			fmt.Printf("Response saved to: %s\n", saveResponseFile)
		} else {
			fmt.Println(res)
		}
	},
}

func getApiResponseImage(args []string) string {
	userArgs := strings.Join(args[0:], " ")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	CheckNilError(err)
	defer client.Close()

	currentGenaiModel := GetConfig("genai_model")
	model := client.GenerativeModel(currentGenaiModel)

	imgData, err := os.ReadFile(imageFilePath)
	CheckNilError(err)

	// Supports multiple image inputs
	prompt := []genai.Part{
		genai.ImageData(imageFileFormat, imgData),
		genai.Text(fmt.Sprintf(userArgs, " in ", respOutputLanguage, " language")),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	CheckNilError(err)

	finalResponse := resp.Candidates[0].Content.Parts[0]
	return fmt.Sprint(finalResponse)

}

func init() {
	imageCmd.Flags().StringVarP(&imageFilePath, "path", "p", "", "Enter the image path")
	imageCmd.Flags().StringVarP(&imageFileFormat, "format", "f", "jpeg", "Enter the image format (jpeg, png, etc.)")
	imageCmd.Flags().StringVarP(&respOutputLanguage, "language", "l", "english", "Enter the language for the output")
	imageCmd.Flags().Float32VarP(&modelTemp, "temperature", "t", 0.5, "Response creativity (0.0-1.0)")
	imageCmd.Flags().BoolVarP(&saveResponse, "save", "s", false, "Save the output to a file")
	imageCmd.Flags().StringVarP(&saveResponseFile, "output", "o", "output.txt", "Output file name")
	errPathF := imageCmd.MarkFlagRequired("path")
	CheckNilError(errPathF)
}
