package server

import (
	"encoding/json"
	"github.com/JIIL07/jcloud/internal/server/middleware"
	"net/http"

	"github.com/JIIL07/jcloud/internal/server/private/commandline"
	"github.com/JIIL07/jcloud/internal/storage"
	"github.com/gorilla/mux"
)

func setupRouter(s *storage.Storage) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.StorageMiddleware(s))

	router.HandleFunc("/", RootHandler).Methods(http.MethodGet).Name("root")
	router.HandleFunc("/api/v1/healthcheck", HealthCheckHandler).Methods(http.MethodGet).Name("healthcheck")

	api := router.PathPrefix("/api/v1").Name("api").Subrouter()
	api.HandleFunc("/login", LoginHandler).Methods(http.MethodPost, http.MethodGet).Name("login")
	api.HandleFunc("/logout", LogoutHandler).Methods(http.MethodGet).Name("logout")

	f := api.PathPrefix("/files").Name("files").Subrouter()
	f.Use(middleware.LoginMiddleware)
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
