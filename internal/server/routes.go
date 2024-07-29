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

	router.HandleFunc("/", rootHandler).Methods(http.MethodGet).Name("root")
	router.HandleFunc("/api/v1/healthcheck", healthCheckHandler).Methods(http.MethodGet).Name("healthcheck")

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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Welcome to CloudFiles API`))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}
