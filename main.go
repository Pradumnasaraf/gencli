package main

import (
	"fmt"
	"os"

	"github.com/Pradumnasaraf/gencli/cmd"
)

func main() {

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	if geminiAPIKey == "" {
		fmt.Println("Please set the GEMINI_API_KEY environment variable. Check the README for more information.")
		return
	}

	cmd.SetDefaultConfig()
	cmd.CheckAPIKey()
	cmd.Execute()
}
