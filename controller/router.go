package controller

import (
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/evaluate", HandleEvaluate).Methods("POST")

	return r
}
