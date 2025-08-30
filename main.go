package main

import (
	"fmt"
	"os"

	"github.com/Pradumnasaraf/gencli/cmd"
)

func main() {

	googleGeminiAPIKey := os.Getenv("GOOGLE_API_KEY")

	if googleGeminiAPIKey == "" {
		fmt.Println("Please set the GOOGLE_API_KEY environment variable. Check the README for more information.")
		return
	}

	cmd.SetDefaultConfig()
	cmd.CheckAPIKey()
	cmd.Execute()
}
