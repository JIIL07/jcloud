package server

import (
	"encoding/json"
	"github.com/JIIL07/jcloud/internal/server/middleware"
	"github.com/JIIL07/jcloud/internal/server/private/commandline"
	"github.com/JIIL07/jcloud/internal/server/static"
	"github.com/JIIL07/jcloud/internal/storage"
	jctx "github.com/JIIL07/jcloud/pkg/ctx"
	"github.com/gorilla/mux"
	"net/http"
)

func setupRouter(s *storage.Storage, b *static.Static) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.StorageMiddleware(s))

	router.HandleFunc("/static/{filename}", b.BinaryHandler).Methods(http.MethodGet)

	router.HandleFunc("/", RootHandler).Methods(http.MethodGet).Name("root")
	router.HandleFunc("/api/v1/healthcheck", HealthCheckHandler).Methods(http.MethodGet).Name("healthcheck")

	api := router.PathPrefix("/api/v1").Name("api").Subrouter()
	api.HandleFunc("/login", LoginHandler).Methods(http.MethodPost, http.MethodGet).Name("login")
	api.HandleFunc("/logout", LogoutHandler).Methods(http.MethodGet).Name("logout")

	user := api.PathPrefix("/user").Name("user").Subrouter()
	user.Use(middleware.UserMiddleware)
	user.HandleFunc("/{user}", CurrentUserHandler).Methods(http.MethodGet).Name("current-user")

	f := user.PathPrefix("/files").Name("files").Subrouter()
	f.Use(middleware.UserMiddleware)
	f.HandleFunc("/get", GetFilesHandler).Methods(http.MethodGet).Name("get-files")
	f.HandleFunc("/upload", AddFileHandler).Methods(http.MethodPost).Name("add-file")
	f.HandleFunc("/download", DownloadFileHandler).Methods(http.MethodGet).Name("download-file")
	f.HandleFunc("/delete", DeleteFileHandler).Methods(http.MethodDelete).Name("delete-file")

	private := router.PathPrefix("/private").Subrouter()
	private.HandleFunc("/admin", commandline.AuthHandler).Methods(http.MethodGet)
	private.HandleFunc("/admin/checkadmin", commandline.CheckHandler).Methods(http.MethodGet)
	private.HandleFunc("/sql", commandline.HandleSQLQuery).Methods(http.MethodGet)
	private.HandleFunc("/cmd", commandline.HandleCmdExec).Methods(http.MethodGet)

	return router
}

func CurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
	}

	user := r.Context().Value("user").(string)
	u, err := s.GetByUsername(user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		http.Error(w, "Failed to encode user", http.StatusInternalServerError)
		return
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Welcome to CloudFiles API`)) // nolint:errcheck
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
