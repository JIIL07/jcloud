package handlers

import (
	"encoding/json"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`Welcome to CloudFiles API`))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
