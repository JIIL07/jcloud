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

	router.HandleFunc("/", rootHandler).Methods("GET").Name("root")
	router.HandleFunc("/api/v1/healthcheck", healthCheckHandler).Methods("GET").Name("healthcheck")

	private := router.PathPrefix("/private").Subrouter()
	private.HandleFunc("/tools/sql", func(w http.ResponseWriter, r *http.Request) {
		ctx := jctx.WithContext(r.Context(), "storage", s)
		r = r.WithContext(ctx)
		commandline.HandleSQLQuery(w, r)
	}).Methods("POST")

	return router
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Welcome to CloudFiles API`))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}
