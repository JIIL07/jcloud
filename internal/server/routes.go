package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler).Methods("GET").Name("root")
	router.HandleFunc("/api/v1/healthcheck", healthCheckHandler).Methods("GET").Name("healthcheck")

	return router
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Welcome to CloudFiles API`))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}
