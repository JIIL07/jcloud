package commandline

import (
	"encoding/json"
	"net/http"
	"os"

	jtoken "github.com/JIIL07/cloudFiles-manager/internal/lib/token"
	"github.com/JIIL07/cloudFiles-manager/internal/storage"
	"github.com/gorilla/sessions"
)

var (
	sessionToken = jtoken.Generate(16)
	store        = sessions.NewCookieStore([]byte(sessionToken))
)

func init() {
	store.Options = &sessions.Options{
		MaxAge:   86400, //24 hours
		Secure:   false,
		HttpOnly: true,
	}
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only method GET is allowed"))
		return
	}

	a := r.URL.Query().Get("admin")

	var u storage.AboutUser

	d := os.Getenv("ADMIN_USER")
	err := json.Unmarshal([]byte(d), &u)
	if err != nil {
		http.Error(w, "Invalid admin user configuration", http.StatusInternalServerError)
		return
	}

	if a == u.Username {
		session, err := store.Get(r, "admin")
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

		w.Write([]byte("Session established"))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Unauthorized"))
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	s, err := store.Get(r, "admin")
	if err != nil {
		respondWithError(w, err)
		return
	}

	if s.IsNew {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	w.Write([]byte("Authorized"))
}
