package commandline

import (
	"encoding/json"
	"io"
	"net/http"

	jctx "github.com/JIIL07/cloudFiles-manager/internal/lib/ctx"
	"github.com/JIIL07/cloudFiles-manager/internal/storage"
)

type Request struct {
	Command string `json:"command"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

var s *storage.Storage

func HandleSQLQuery(w http.ResponseWriter, r *http.Request) {
	var ok bool
	s, ok = jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	rows, err := s.Query(req.Command)
	if err != nil {
		respondWithError(w, err)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	results, err = storage.ParseRows(rows)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, results)
}

func setAdmin() string {
	err := s.Admin()
	if err != nil {
		return err.Error()
	}
	return "Admin user set"
}

func respondWithError(w http.ResponseWriter, err error) {
	respondWithJSON(w, Response{Error: err.Error()})
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
