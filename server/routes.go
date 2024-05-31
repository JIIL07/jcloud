package server

import (
	"encoding/json"
	"net/http"

	cloudfiles "github.com/JIIL07/cloudFiles-manager/client"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes(s *ServerContext) *httprouter.Router {
	router := httprouter.New()

	// Define the available routes
	router.Handle(http.MethodGet, "/", HandlerAdapter(app.TextHandler))
	router.Handle(http.MethodGet, "/files", HandlerAdapter(app.SetFilesHandler(s)))
	router.Handle(http.MethodPost, "/addfiles", HandlerAdapter(app.AddFileHandler(s)))
	router.Handle(http.MethodGet, "/deletedfiles", HandlerAdapter(app.SetDeletedFilesHandler(s)))

	return router
}

func HandlerAdapter(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h(w, r)
	}
}

func (app *application) TextHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`
Welcome to the cloudfiles API. You can use the following endpoints:
GET /files
GET /deletedfiles
POST /addfiles
POST /deletefiles
POST /updatefiles
	`))

	app.logger.Printf("Server detected / entering")
}

func (app *application) SetFilesHandler(s *ServerContext) http.HandlerFunc {
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

func (app *application) SetDeletedFilesHandler(s *ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := s.Ctx.List("deletedfiles")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func (app *application) AddFileHandler(s *ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fileContext cloudfiles.FileContext
		if err := json.NewDecoder(r.Body).Decode(&fileContext); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			app.errlogger.Printf("error decoding: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		app.logger.Printf("Successfully written file")
	}
}
