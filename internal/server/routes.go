package server

import (
	"encoding/json"
	"net/http"

	jctx "github.com/JIIL07/cloudFiles-manager/internal/lib/ctx"
	"github.com/JIIL07/cloudFiles-manager/internal/server/private/commandline"
	"github.com/JIIL07/cloudFiles-manager/internal/storage"
	"github.com/gorilla/mux"
)

func setupRouter(s *storage.Storage) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", RootHandler).Methods(http.MethodGet).Name("root")
	router.HandleFunc("/api/v1/healthcheck", HealthCheckHandler).Methods(http.MethodGet).Name("healthcheck")

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		ctx := jctx.WithContext(r.Context(), "storage", s)
		r = r.WithContext(ctx)
		LoginHandler(w, r)
	}).Methods(http.MethodPost)
	api.HandleFunc("/login/check", func(w http.ResponseWriter, r *http.Request) {
		ctx := jctx.WithContext(r.Context(), "storage", s)
		r = r.WithContext(ctx)
		LoginCheckHandler(w, r)
	}).Methods(http.MethodGet)

	f := api.PathPrefix("/files").Subrouter()
	f.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		ctx := jctx.WithContext(r.Context(), "storage", s)
		r = r.WithContext(ctx)
		GetFilesHandler(w, r)
	}).Methods(http.MethodGet)
	f.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

	})

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
	w.Write([]byte(`Welcome to CloudFiles API`))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}
