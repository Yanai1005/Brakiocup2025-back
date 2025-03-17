package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".envファイルが見つからないか読み込めません")
	}

	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEYが必要です")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	prompt := genai.Text("日本の現在の総理大臣は誰ですか？")

	resp, err := model.GenerateContent(ctx, prompt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("===answer===")
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			fmt.Println(part)
		}
	}

	fmt.Println("============")
}
