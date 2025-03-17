package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Brakiocup2025-back/model"
)

func HandleEvaluate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")

	log.Printf("リクエストを受信: %s %s", r.Method, r.URL.Path)
	log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("リクエスト読み込みエラー: %v", err)
		sendErrorResponse(w, "リクエストの読み込みに失敗", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	log.Printf("リクエスト本文: %s", string(body))
	var evalReq model.EvaluationRequest
	if err := json.Unmarshal(body, &evalReq); err != nil {
		log.Printf("JSONパースエラー: %v", err)
		sendErrorResponse(w, "JSONのパースに失敗", http.StatusBadRequest)
		return
	}

	if evalReq.Content == "" {
		sendErrorResponse(w, "評価するコンテンツが必要", http.StatusBadRequest)
		return
	}

	scores, err := model.EvaluateReadme(evalReq.Content)
	if err != nil {
		sendErrorResponse(w, "評価に失敗: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(scores)
}

func sendErrorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
	evalResp := model.EvaluationResponse{
		Error: errorMsg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(evalResp)
}
