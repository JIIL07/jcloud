package server

import (
	"crypto"
	"encoding/json"
	"github.com/JIIL07/cloudFiles-manager/internal/lib/cookies"
	jctx "github.com/JIIL07/cloudFiles-manager/internal/lib/ctx"
	jhash "github.com/JIIL07/cloudFiles-manager/internal/lib/hash"
	"github.com/JIIL07/cloudFiles-manager/internal/lib/ip"
	"github.com/JIIL07/cloudFiles-manager/internal/lib/role"
	"github.com/JIIL07/cloudFiles-manager/internal/storage"
	"github.com/gorilla/sessions"
	"net/http"
)

var credentials *CurrentUser

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}

	var credentials CurrentUser
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

	credentials.UserData.Password = jhash.Hash(credentials.UserData.Password)
	credentials.UserData.Protocol = crypto.SHA256.String()
	f := storage.Admin(&credentials.UserData)
	credentials.Role = role.Set(f)
	credentials.UserData.Admin = credentials.Role
	credentials.NetworkDetails.IP = ip.GetIPAddress(r)

	if err = s.SaveNewUser(&credentials.UserData); err != nil {
		http.Error(w, "Failed to save new user", http.StatusInternalServerError)
		return
	}

	session, err := cookies.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
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

	w.Write([]byte("Authorization successful"))
}

func LoginCheckHandler(w http.ResponseWriter, r *http.Request) {
	s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}
	e, err := s.CheckExistence(credentials.UserData.Username)

	if err != nil {
		http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
		return
	}

	if !e {
		http.Error(w, "User doesn't exist", http.StatusBadRequest)
		return
	}
	w.Write([]byte("Authorized"))
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
	w.Write([]byte("Session cleared"))
}
