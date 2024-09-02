package commandline

import (
	jjson "github.com/JIIL07/jcloud/pkg/json"
	"net/http"

	"github.com/JIIL07/jcloud/internal/storage"
	"github.com/JIIL07/jcloud/pkg/cookies"
	"github.com/JIIL07/jcloud/pkg/ctx"
	"github.com/JIIL07/jcloud/pkg/parsers"
)

var s *storage.Storage

func HandleSQLQuery(w http.ResponseWriter, r *http.Request) {
	var ok bool
	s, ok = jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}

	store, err := cookies.Store.Get(r, "admin")
	if err != nil {
		jjson.RespondWithError(w, err)
		return
	}

	if store.IsNew {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized")) // nolint:errcheck
		return
	}

	var req jjson.Request
	req.Command = r.URL.Query().Get("command")

	rows, err := s.Query(req.Command)
	if err != nil {
		jjson.RespondWithError(w, err)
		return
	}
	defer rows.Close() // nolint:errcheck

	var results []map[string]interface{}
	results, err = parsers.ParseRows(rows)
	if err != nil {
		jjson.RespondWithError(w, err)
		return
	}

	jjson.RespondWithJSON(w, results)
}
