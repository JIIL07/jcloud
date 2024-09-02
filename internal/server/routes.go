package server

import (
	"encoding/json"
	"net/http"

	"github.com/JIIL07/jcloud/internal/server/private/commandline"
	"github.com/JIIL07/jcloud/internal/storage"
	"github.com/JIIL07/jcloud/pkg/ctx"
	"github.com/gorilla/mux"
)

func setupRouter(s *storage.Storage) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", RootHandler).Methods(http.MethodGet).Name("root")
	router.HandleFunc("/api/v1/healthcheck", HealthCheckHandler).Methods(http.MethodGet).Name("healthcheck")

	api := router.PathPrefix("/api/v1").Name("api").Subrouter()
	api.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		ctx := jctx.WithContext(r.Context(), "storage", s)
		r = r.WithContext(ctx)
		LoginHandler(w, r)
	}).Methods(http.MethodPost, http.MethodGet).Name("login")
	api.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		ctx := jctx.WithContext(r.Context(), "storage", s)
		r = r.WithContext(ctx)
		LogoutHandler(w, r)
	}).Methods(http.MethodGet).Name("logout")

	f := api.PathPrefix("/files").Name("files").Subrouter()
	f.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		ctx := jctx.WithContext(r.Context(), "storage", s)
		r = r.WithContext(ctx)
		GetFilesHandler(w, r)
	}).Methods(http.MethodGet).Name("get-files")
	f.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		ctx := jctx.WithContext(r.Context(), "storage", s)
		r = r.WithContext(ctx)
		AddFileHandler(w, r)
	}).Methods(http.MethodPost).Name("add-file")
	//f.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
	//	ctx := jctx.WithContext(r.Context(), "storage", s)
	//	r = r.WithContext(ctx)
	//	DownloadFileHandler(w, r)
	//}).Methods(http.MethodGet).Name("download-file")

	f.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		ctx := jctx.WithContext(r.Context(), "storage", s)
		r = r.WithContext(ctx)
		DeleteFileHandler(w, r)
	}).Methods(http.MethodDelete).Name("delete-file")

	private := router.PathPrefix("/private").Subrouter()
	private.HandleFunc("/admin", commandline.AuthHandler).Methods(http.MethodGet)
	private.HandleFunc("/admin/checkadmin", commandline.CheckHandler).Methods(http.MethodGet)
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
	w.Write([]byte(`Welcome to CloudFiles API`)) // nolint:errcheck
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
