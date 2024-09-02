package json

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

func RespondWithError(w http.ResponseWriter, err error) {
	RespondWithJSON(w, Response{Error: err.Error()})
}

func RespondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		RespondWithError(w, err)
	}
}
