package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Brakiocup2025-back/model"
)

func HandleUserAnalysis(w http.ResponseWriter, r *http.Request) {
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
		sendUserAnalysisErrorResponse(w, "リクエストの読み込みに失敗", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("リクエスト本文: %s", string(body))

	var analysisReq model.UserAnalysisRequest
	if err := json.Unmarshal(body, &analysisReq); err != nil {
		log.Printf("JSONパースエラー: %v", err)
		sendUserAnalysisErrorResponse(w, "JSONのパースに失敗", http.StatusBadRequest)
		return
	}

	if analysisReq.Username == "" {
		sendUserAnalysisErrorResponse(w, "ユーザー名が必要です", http.StatusBadRequest)
		return
	}

	//最大10にする
	if analysisReq.MaxRepos <= 0 || analysisReq.MaxRepos > 10 {
		analysisReq.MaxRepos = 10
	}

	log.Printf("ユーザー分析開始: %s (最大%d件)", analysisReq.Username, analysisReq.MaxRepos)

	analysisResp, err := model.AnalyzeUserRepositories(analysisReq.Username, analysisReq.MaxRepos)
	if err != nil {
		log.Printf("ユーザー分析エラー: %v", err)
		sendUserAnalysisErrorResponse(w, "ユーザーの分析に失敗: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("ユーザー分析完了: %s (%d/%d リポジトリ分析)",
		analysisReq.Username,
		analysisResp.AnalyzedCount,
		analysisResp.RepositoryCount)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(analysisResp)
}

func sendUserAnalysisErrorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
	response := model.UserAnalysisResponse{
		Error: errorMsg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
