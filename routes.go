package cloudfiles

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes(s *ServerContext) *httprouter.Router {
	router := httprouter.New()

	// Define the available routes
	router.Handle(http.MethodGet, "/", HandlerAdapter(TextHandler))
	router.Handle(http.MethodGet, "/files", HandlerAdapter(SetFilesHandler(s)))
	router.Handle(http.MethodGet, "/deletedfiles", HandlerAdapter(SetDeletedFilesHandler(s)))

	return router
}

func HandlerAdapter(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h(w, r)
	}
}

func TextHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`Welcome to the cloudfiles API. You can use the following endpoints:
GET /files
GET /deletedfiles`))
}

func SetFilesHandler(s *ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := s.Ctx.List("files")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func SetDeletedFilesHandler(s *ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := s.Ctx.List("files")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func (app *application) SetAddFileHandler(ctx *FileContext) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var fileContext FileContext
		if err := json.NewDecoder(r.Body).Decode(&fileContext); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := ctx.Add(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
