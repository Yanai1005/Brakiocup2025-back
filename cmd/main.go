package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Brakiocup2025-back/controller"
	"github.com/Brakiocup2025-back/model"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".envファイルが見つからないか読み込めません")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEYが必要です")
	}

	err = model.InitGeminiClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer model.CloseGeminiClient()

	router := controller.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("サーバーを起動 ポート: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
