package server

import (
	"crypto"
	"encoding/json"
	"net/http"

	jctx "github.com/JIIL07/cloudFiles-manager/internal/lib/ctx"
	jhash "github.com/JIIL07/cloudFiles-manager/internal/lib/hash"
	"github.com/JIIL07/cloudFiles-manager/internal/lib/ip"
	"github.com/JIIL07/cloudFiles-manager/internal/lib/role"
	"github.com/JIIL07/cloudFiles-manager/internal/server/private/commandline"
	"github.com/JIIL07/cloudFiles-manager/internal/storage"
	"github.com/gorilla/mux"
)

func setupRouter(s *storage.Storage) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", RootHandler).Methods(http.MethodGet).Name("root")
	router.HandleFunc("/api/v1/healthcheck", HealthCheckHandler).Methods(http.MethodGet).Name("healthcheck")

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		ctx := jctx.WithContext(r.Context(), "storage", s)
		r = r.WithContext(ctx)
		AuthorizationHandler(w, r)
	}).Methods(http.MethodPost)

	private := router.PathPrefix("/private").Subrouter()
	private.HandleFunc("/auth", commandline.AuthHandler).Methods(http.MethodGet)
	private.HandleFunc("/auth/checkadmin", commandline.CheckHandler).Methods(http.MethodGet)
	private.HandleFunc("/sql", func(w http.ResponseWriter, r *http.Request) {
		ctx := jctx.WithContext(r.Context(), "storage", s)
		r = r.WithContext(ctx)
		commandline.HandleSQLQuery(w, r)
	}).Methods(http.MethodGet)

	private.HandleFunc("/cmd", func(w http.ResponseWriter, r *http.Request) {
		ctx := jctx.WithContext(r.Context(), "storage", s)
		r = r.WithContext(ctx)
		commandline.HandleCmdExec(w, r)
	}).Methods(http.MethodGet)

	return router
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Welcome to CloudFiles API`))
}

func AuthorizationHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Write([]byte("Authorization successful"))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}
