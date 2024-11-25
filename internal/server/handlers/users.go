package handlers

import (
	"crypto"
	"encoding/json"
	"fmt"
	"github.com/JIIL07/jcloud/internal/server/cookies"
	"github.com/JIIL07/jcloud/internal/server/storage"
	"github.com/JIIL07/jcloud/internal/server/utils"
	"github.com/JIIL07/jcloud/pkg/ip"
	"net/http"
)

type CurrentUser struct {
	UserData storage.User
	Role     int
	Network  Connection
}

type Connection struct {
	IP string
}

func Login(w http.ResponseWriter, r *http.Request, c *CurrentUser) error {
	session := cookies.GetSession(r, "user-session")

	session.Values["user"] = true
	session.Values["role"] = c.Role
	session.Values["username"] = c.UserData.Username

	err := cookies.Store.Save(r, w, session)
	if err != nil {
		return err
	}

	return nil
}

func SaveUser(w http.ResponseWriter, r *http.Request) *CurrentUser {
	var user CurrentUser

	s := utils.ProvideStorage(r, w)

	err := json.NewDecoder(r.Body).Decode(&user.UserData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return nil
	}

	e, err := s.CheckUser(user.UserData.Username)
	if err != nil {
		http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
		return nil
	}

	if e {
		http.Error(w, "User exists", http.StatusBadRequest)
		return nil
	}

	user.UserData.Protocol = crypto.SHA256.String()
	user.UserData.Admin = storage.Admin(&user.UserData).Int()
	user.Network.IP = ip.GetIPAddress(r)

	if err = s.SaveNewUser(&user.UserData); err != nil {
		http.Error(w, "Failed to save new user", http.StatusInternalServerError)
		return nil
	}

	return &user
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := SaveUser(w, r)
	err := Login(w, r, user)
	if err != nil {
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Authorization successful")) // nolint:errcheck

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookies.ClearSession(w, r)
	w.Write([]byte("Session cleared")) // nolint:errcheck
}

func CurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode user", http.StatusInternalServerError)
		return
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	w.Write([]byte(fmt.Sprintf("Current user: %v\n", user.Username))) // nolint:errcheck

	s := utils.ProvideStorage(r, w)
	files, err := s.GetAllFiles(user.UserID)
	if err != nil {
		http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
		return
	}

	var t int
	for i := range files {
		t += 1
		json.NewEncoder(w).Encode(files[i].Metadata.Name + "." + files[i].Metadata.Extension) // nolint:errcheck
	}

	w.Write([]byte(fmt.Sprintf("\ntotal files: %d", t))) // nolint:errcheck
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	u := utils.ProvideUser(r, w)
	s := utils.ProvideStorage(r, w)
	err := s.DeleteUser(u.Username)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted")) // nolint:errcheck
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	u := utils.ProvideUser(r, w)
	s := utils.ProvideStorage(r, w)
	newPassword := r.URL.Query().Get("password")
	newEmail := r.URL.Query().Get("email")
	updates := make(map[string]interface{})

	if newPassword == "" && newEmail == "" {
		http.Error(w, "No updates provided", http.StatusBadRequest)
		return
	}

	if newPassword != "" {
		updates["password"] = newPassword
	}

	if newEmail != "" {
		updates["email"] = newEmail
	}

	err := s.UpdateUserInfo(u.Username, updates)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated")) // nolint:errcheck
}
