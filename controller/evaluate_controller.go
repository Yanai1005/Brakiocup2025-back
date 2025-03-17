package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Brakiocup2025-back/model"
)

func HandleEvaluate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, "リクエストの読み込みに失敗", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var evalReq model.EvaluationRequest
	if err := json.Unmarshal(body, &evalReq); err != nil {
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
