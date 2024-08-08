package server

import (
	"encoding/json"
	"github.com/JIIL07/cloudFiles-manager/internal/lib/cookies"
	jctx "github.com/JIIL07/cloudFiles-manager/internal/lib/ctx"
	"github.com/JIIL07/cloudFiles-manager/internal/storage"
	"net/http"
)

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}
	session, err := cookies.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	u, err := s.GetByUsername(session.Values["username"].(string))
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}
	files, err := s.GetAllFiles(u)
	if err != nil {
		http.Error(w, "Failed to get files", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}
