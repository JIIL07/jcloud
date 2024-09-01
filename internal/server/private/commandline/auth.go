package commandline

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/JIIL07/jcloud/internal/lib/cookies"
	"github.com/JIIL07/jcloud/internal/storage"
	"github.com/gorilla/sessions"
)

var u storage.UserData

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only method GET is allowed")) // nolint:errcheck
		return
	}

	a := r.URL.Query().Get("admin")

	d := os.Getenv("ADMIN_USER")
	err := json.Unmarshal([]byte(d), &u)
	if err != nil {
		http.Error(w, "Invalid admin user configuration", http.StatusInternalServerError)
		return
	}

	if a == u.Username {
		session, err := cookies.Store.Get(r, "admin")
		if err != nil {
			respondWithError(w, err)
			return
		}

		session.Values["admin"] = true
		session.Values["sql"] = true
		session.Values["cmd"] = true

		err = sessions.Save(r, w)
		if err != nil {
			respondWithError(w, err)
			return
		}

		w.Write([]byte("Session established")) // nolint:errcheck
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Unauthorized")) // nolint:errcheck
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	s, err := cookies.Store.Get(r, "admin")
	if err != nil {
		respondWithError(w, err)
		return
	}

	if s.IsNew {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized")) // nolint:errcheck
		return
	}

	w.Write([]byte("Admin authorized")) // nolint:errcheck
}
