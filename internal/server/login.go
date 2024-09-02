package server

import (
	"crypto"
	"encoding/json"
	"github.com/JIIL07/jcloud/internal/storage"
	"github.com/JIIL07/jcloud/pkg/cookies"
	"github.com/JIIL07/jcloud/pkg/ctx"
	"github.com/JIIL07/jcloud/pkg/ip"
	"github.com/gorilla/sessions"
	"net/http"
)

var ActualUser = &CurrentUser{}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials CurrentUser
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

	switch r.Method {
	case http.MethodPost:
		{
			err := json.NewDecoder(r.Body).Decode(&credentials.UserData)
			if err != nil {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}

			e, err := s.CheckExistence(credentials.UserData.Username)
			if err != nil {
				http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
				return
			}

			if e {
				http.Error(w, "User exists", http.StatusBadRequest)
				return
			}

			credentials.UserData.Protocol = crypto.SHA256.String()
			credentials.UserData.Admin = storage.Admin(&credentials.UserData).Int()

			if err = s.SaveNewUser(&credentials.UserData); err != nil {
				http.Error(w, "Failed to save new user", http.StatusInternalServerError)
				return
			}

			session.Values["user"] = true
			session.Values["username"] = credentials.UserData.Username
			session.Values["email"] = credentials.UserData.Email

			err = sessions.Save(r, w)
			if err != nil {
				http.Error(w, "Failed to save session", http.StatusInternalServerError)
				return
			}

			w.Write([]byte("Authorization successful")) // nolint:errcheck
		}
	case http.MethodGet:
		{
			credentials.UserData.Username = r.URL.Query().Get("username")
			e, err := s.CheckExistence(credentials.UserData.Username)
			if err != nil {
				http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
				return
			}

			if !e {
				http.Error(w, "User doesn't exist", http.StatusBadRequest)
				return
			}

			session, err = cookies.Store.Get(r, "user-session")
			if err != nil {
				http.Error(w, "Failed to get session", http.StatusInternalServerError)
				return
			}

			ActualUser.UserData, err = s.GetByUsername(credentials.UserData.Username)
			ActualUser.NetworkDetails.IP = ip.GetIPAddress(r)
			if err != nil {
				http.Error(w, "Failed to get user", http.StatusInternalServerError)
				return
			}

			err = session.Save(r, w)
			if err != nil {
				http.Error(w, "Failed to save session", http.StatusInternalServerError)
				return
			}

			w.Write([]byte("Authorized")) // nolint:errcheck
		}
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := cookies.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	session.Values["user"] = false
	err = sessions.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Session cleared")) // nolint:errcheck
}
