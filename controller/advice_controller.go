package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Brakiocup2025-back/model"
)

func HandleReadmeAdvice(w http.ResponseWriter, r *http.Request) {
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
		sendAdviceErrorResponse(w, "リクエストの読み込みに失敗", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("リクエスト本文: %s", string(body))

	var adviceReq model.ReadmeAdviceRequest
	if err := json.Unmarshal(body, &adviceReq); err != nil {
		log.Printf("JSONパースエラー: %v", err)
		sendAdviceErrorResponse(w, "JSONのパースに失敗", http.StatusBadRequest)
		return
	}

	var content string
	var fetchErr error

	// コンテンツが直接提供されているか、GitHubリポジトリから取得するか
	if adviceReq.Content != "" {
		content = adviceReq.Content
	} else if adviceReq.Owner != "" && adviceReq.Repo != "" {
		// GitHubからREADMEを取得
		content, fetchErr = model.FetchReadme(adviceReq.Owner, adviceReq.Repo)
		if fetchErr != nil {
			log.Printf("README取得エラー: %v", fetchErr)
			sendAdviceErrorResponse(w, "GitHubからのREADME取得に失敗: "+fetchErr.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		sendAdviceErrorResponse(w, "コンテンツか、オーナーとリポジトリ名のどちらかが必要です", http.StatusBadRequest)
		return
	}

	// READMEを評価してアドバイスを生成
	adviceResp, err := model.GenerateReadmeAdviceWithImprovement(content)
	if err != nil {
		log.Printf("アドバイス生成エラー: %v", err)
		sendAdviceErrorResponse(w, "アドバイスの生成に失敗: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(adviceResp)
}

func sendAdviceErrorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
	response := model.ReadmeAdviceResponse{
		Error: errorMsg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
