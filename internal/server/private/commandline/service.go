package commandline

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	Command string
	Args    string
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func respondWithError(w http.ResponseWriter, err error) {
	respondWithJSON(w, Response{Error: err.Error()})
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
