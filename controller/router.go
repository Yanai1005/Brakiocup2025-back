package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(corsMiddleware)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok","message":"API server is running"}`)
	}).Methods("GET", "OPTIONS")
	r.HandleFunc("/evaluate", HandleEvaluate).Methods("POST", "OPTIONS")
	r.HandleFunc("/github/readme", HandleGitHubReadme).Methods("POST", "OPTIONS")
	r.HandleFunc("/github/advice", HandleReadmeAdvice).Methods("POST", "OPTIONS")
	r.HandleFunc("/github/analysis", HandleUserAnalysis).Methods("POST", "OPTIONS")
	return r
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
