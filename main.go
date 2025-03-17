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

	if len(os.Args) < 2 {
		log.Fatal("使用方法: go run main.go <TEXT_FILE_PATH>")
	}
	textFilePath := os.Args[1]

	textContent, err := os.ReadFile(textFilePath)
	if err != nil {
		log.Fatalf("テキストファイルの読み込みに失敗しました: %v", err)
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

	evaluationPrompt := `以下のREADMEファイルを評価してください。評価基準は次のとおりです：

1. 明確さ: プロジェクトの目的を明確に説明されているか（0〜10点）
2. 完全性:  環境構築、インストール手順、使用方法などの必要な情報がすべて含まれているか（0〜10点）
3. 構造化: 情報が論理的に整理されているか（0〜10点）
4. 例示: 使用例やコード例が含まれているか（0〜10点）
5. 可読性: 文章が読みやすいか（0〜10点）

各基準の点数を示してください。

評価対象のテキスト:
"""
` + string(textContent) + `
"""

合計点数を提供してください。`

	resp, err := model.GenerateContent(ctx, genai.Text(evaluationPrompt))
	if err != nil {
		log.Fatalf("生成に失敗しました: %v", err)
	}

	fmt.Println("=== テキスト評価結果 ===")
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			fmt.Println(part)
		}
	}
	fmt.Println("=======================")
}
