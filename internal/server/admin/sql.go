package admin

import (
	"net/http"

	"github.com/JIIL07/jcloud/internal/server/cookies"
	"github.com/JIIL07/jcloud/internal/server/utils"
	jjson "github.com/JIIL07/jcloud/pkg/json"

	"github.com/JIIL07/jcloud/pkg/parsers"
)

func HandleSQLQuery(w http.ResponseWriter, r *http.Request) {
	s := utils.ProvideStorage(r, w)
	store := cookies.GetSession(r, "admin")

	if store.IsNew {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized"))
		return
	}

	var req jjson.Request
	req.Command = r.URL.Query().Get("command")

	rows, err := s.Query(req.Command)
	if err != nil {
		jjson.RespondWithError(w, err)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	results, err = parsers.ParseRows(rows)
	if err != nil {
		jjson.RespondWithError(w, err)
		return
	}

	jjson.RespondWithJSON(w, results)
}
