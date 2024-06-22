package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Config() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(".env file not found. Loading the system environment variables.")
	}
}
