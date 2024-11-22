package middleware

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/server/cookies"
	"github.com/JIIL07/jcloud/internal/server/storage"
	"github.com/JIIL07/jcloud/pkg/ctx"
	"net/http"
)

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
		if !ok {
			http.Error(w, "Storage not found", http.StatusInternalServerError)
			return
		}
		session, err := cookies.Store.Get(r, "user-session")
		if err != nil || session.IsNew {
			cookies.ClearSession(w, r)
			http.Redirect(w, r, "api/v1/login", http.StatusSeeOther)
			return
		}

		if session.Values["username"] == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		u, err := s.GetByUsername(session.Values["username"].(string))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get user: %v", err), http.StatusInternalServerError)
			return
		}

		ctx := jctx.WithContext(r.Context(), "user", &u)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func StorageMiddleware(s *storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := jctx.WithContext(r.Context(), "storage", s)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
