package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	configFileDir  string = ".gencli"
	configFileName string = "config"
	configFileType string = "yaml"
	defaultModel   string = "gemini-2.5-pro"
)

func SetDefaultConfig() {
	homeDir := getHomeDir()

	if _, err := os.Stat(homeDir + "/" + configFileDir); os.IsNotExist(err) {
		if err := os.Mkdir(homeDir+"/"+configFileDir, 0755); err != nil {
			log.Fatalf("Error creating config directory: %v", err)
		}

		configFilePath := homeDir + "/" + configFileDir
		viper.SetConfigName(configFileName)
		viper.SetConfigType(configFileType)
		viper.AddConfigPath(configFilePath)

		viper.Set("genai_model", defaultModel)
		if err := viper.WriteConfigAs(configFilePath + "/" + configFileName + "." + configFileType); err != nil {
			log.Fatalf("Error writing config file: %v", err)
		}
		return
	}
}

func UpdateConfig(key string, value string) {
	homeDir := getHomeDir()

	configFilePath := homeDir + "/" + configFileDir
	viper.Set(key, value)

	if err := viper.WriteConfigAs(configFilePath + "/" + configFileName + "." + configFileType); err != nil {
		log.Fatalf("Error writing config file: %v", err)
	}
}

func GetConfig(key string) string {
	// Ensure Viper has the correct config file set
	if viper.ConfigFileUsed() == "" {
		homeDir := getHomeDir()
		configFilePath := homeDir + "/" + configFileDir + "/" + configFileName + "." + configFileType
		viper.SetConfigFile(configFilePath)
		viper.SetConfigType(configFileType)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}
	}

	return viper.GetString(key)
}

func CheckAPIKey() {
	googleGeminiAPIKey := os.Getenv("GOOGLE_API_KEY")

	if googleGeminiAPIKey == "" {
		fmt.Println("Please set the GOOGLE_API_KEY environment variable. Check the README for more information.")
		return
	}
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to get user home directory to create config file: %v", err)
	}
	return homeDir
}

func CheckNilError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
