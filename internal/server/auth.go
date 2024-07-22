package server

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func loginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Userame, user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		session, _ := store.Get(r, "session")
		session.Values["authenticated"] = true
		session.Values["user"] = username

		if username == "JIIL" && password == "juice" {
			session.Values["role"] = "admin"
		} else {
			session.Values["role"] = "user"
		}

		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400, // 24 hours
			HttpOnly: true,
		}

		if err := session.Save(r, w); err != nil {
			http.Error(w, "Failed to save session", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	role, _ := session.Values["role"].(string)

	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Write([]byte(`Welcome Admin!`))
}
