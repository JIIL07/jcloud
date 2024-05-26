package cloudfiles

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type ServerContext struct {
	Ctx *FileContext
}

func NewServerContext(db *sql.DB) *ServerContext {
	ctx := &FileContext{
		DB:   db,
		Info: &Info{},
	}
	return &ServerContext{
		Ctx: ctx,
	}
}

func (s *ServerContext) SetFilesHandler(ctx *FileContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := ctx.List("files")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func (s *ServerContext) SetDeletedFilesHandler(ctx *FileContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := ctx.List("files")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func (s *ServerContext) Start() error {
	go func() error {
		http.HandleFunc("/files", s.SetFilesHandler(s.Ctx))
		http.HandleFunc("/deletedfiles", s.SetDeletedFilesHandler(s.Ctx))
		err := http.ListenAndServe(":8080", nil)
		return err
	}()
	return nil
}
