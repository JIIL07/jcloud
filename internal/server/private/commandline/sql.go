package commandline

import (
	"net/http"

	jctx "github.com/JIIL07/cloudFiles-manager/internal/lib/ctx"
	"github.com/JIIL07/cloudFiles-manager/internal/storage"
)

var s *storage.Storage

func HandleSQLQuery(w http.ResponseWriter, r *http.Request) {
	var ok bool
	s, ok = jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}

	var req Request
	req.Command = r.FormValue("command")
	req.Token = r.FormValue("token")

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
