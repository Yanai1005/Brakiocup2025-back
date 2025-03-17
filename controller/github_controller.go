package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Brakiocup2025-back/model"
)

func HandleGitHubReadme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")

	log.Printf("リクエストを受信: %s %s", r.Method, r.URL.Path)
	log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("リクエスト読み込みエラー: %v", err)
		sendGitHubErrorResponse(w, "リクエストの読み込みに失敗", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("リクエスト本文: %s", string(body))

	var githubReq model.GitHubRequest
	if err := json.Unmarshal(body, &githubReq); err != nil {
		log.Printf("JSONパースエラー: %v", err)
		sendGitHubErrorResponse(w, "JSONのパースに失敗", http.StatusBadRequest)
		return
	}

	if githubReq.Owner == "" || githubReq.Repo == "" {
		sendGitHubErrorResponse(w, "オーナーとリポジトリ名が必要です", http.StatusBadRequest)
		return
	}

	content, err := model.FetchReadme(githubReq.Owner, githubReq.Repo)
	if err != nil {
		log.Printf("README取得エラー: %v", err)
		sendGitHubErrorResponse(w, "GitHubからのREADME取得に失敗: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := model.GitHubResponse{
		Content: content,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func sendGitHubErrorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
	response := model.GitHubResponse{
		Error: errorMsg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
