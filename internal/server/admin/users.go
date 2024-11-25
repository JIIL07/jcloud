package admin

import (
	"encoding/json"
	"github.com/JIIL07/jcloud/internal/server/utils"
	"net/http"
)

func AllUsersHandler(w http.ResponseWriter, r *http.Request) {
	s := utils.ProvideStorage(r, w)
	users, err := s.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
}
