package model

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	geminiClient *genai.Client
	geminiModel  *genai.GenerativeModel
)

func InitGeminiClient(apiKey string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return err
	}

	geminiClient = client
	geminiModel = client.GenerativeModel("gemini-1.5-flash")

	return nil
}

func CloseGeminiClient() {
	if geminiClient != nil {
		geminiClient.Close()
	}
}

func GenerateGeminiResponse(prompt string) (string, error) {
	ctx := context.Background()

	resp, err := geminiModel.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	var result string
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			result += string(part.(genai.Text))
		}
	}

	return result, nil
}
