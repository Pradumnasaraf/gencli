package cmd

import (
	"context"
	"os"
	"strings"

	"google.golang.org/genai"
)

// ModelInfo represents a selectable GenAI model.
type ModelInfo struct {
	DisplayName string
	ModelID     string
}

// defaultModels is the fallback list used when remote fetching fails.
// This is the single source of truth for predefined models.
var defaultModels = []ModelInfo{
	{DisplayName: "Gemini 3 Pro", ModelID: "gemini-3-pro-preview"},
	{DisplayName: "Gemini 3 Flash", ModelID: "gemini-3-flash-preview"},
	{DisplayName: "Gemini 3.5 Flash", ModelID: "gemini-3.5-flash"},
	{DisplayName: "Gemini 2.5 Pro", ModelID: "gemini-2.5-pro"},
	{DisplayName: "Gemini 2.5 Flash", ModelID: "gemini-2.5-flash"},
	{DisplayName: "Gemini 2.5 Flash-Lite", ModelID: "gemini-2.5-flash-lite"},
	{DisplayName: "Gemini 2.0 Flash", ModelID: "gemini-2.0-flash"},
	{DisplayName: "Gemini 2.0 Flash-Lite", ModelID: "gemini-2.0-flash-lite"},
}

// getModels fetches available Gemini models from the API.
// Falls back to predefined models on any error (network, API, empty response).
func getModels(ctx context.Context) []ModelInfo {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GOOGLE_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return defaultModels
	}

	var models []ModelInfo
	for m, err := range client.Models.All(ctx) {
		if err != nil {
			return defaultModels
		}
		if strings.HasPrefix(m.Name, "gemini-") {
			models = append(models, ModelInfo{
				DisplayName: m.DisplayName,
				ModelID:     m.Name,
			})
		}
	}

	if len(models) == 0 {
		return defaultModels
	}

	return models
}