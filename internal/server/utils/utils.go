package utils

import (
	"github.com/JIIL07/jcloud/internal/server/storage"
	jctx "github.com/JIIL07/jcloud/pkg/ctx"
	"net/http"
)

func ProvideStorage(r *http.Request, w http.ResponseWriter) *storage.Storage {
	s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return nil
	}
	return s
}

func ProvideUser(r *http.Request, w http.ResponseWriter) *storage.User {
	u, ok := jctx.FromContext[*storage.User](r.Context(), "user")
	if !ok {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return nil
	}
	return u
}
