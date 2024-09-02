package server

import (
	"encoding/json"
	"github.com/JIIL07/jcloud/internal/storage"
	"github.com/JIIL07/jcloud/pkg/cookies"
	jctx "github.com/JIIL07/jcloud/pkg/ctx"
	"net/http"
)

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}
	session, err := cookies.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	if session.Values["username"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	u, err := s.GetByUsername(session.Values["username"].(string))
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	f := &storage.File{UserID: u.UserID}
	files, err := s.GetAllFiles(f)
	if err != nil {
		http.Error(w, "Failed to get files", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(files)
	if err != nil {
		http.Error(w, "Failed to encode files", http.StatusInternalServerError)
		return
	}
}

func AddFileHandler(w http.ResponseWriter, r *http.Request) {
	s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}
	session, err := cookies.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	if session.Values["username"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	u, err := s.GetByUsername(session.Values["username"].(string))
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}
	var files []storage.File
	err = json.NewDecoder(r.Body).Decode(&files)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, file := range files {
		file.UserID = u.UserID
		err = s.AddFile(&file)
		if err != nil {
			http.Error(w, "Failed to add file"+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File added")) // nolint:errcheck
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}
	session, err := cookies.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	if session.Values["username"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	u, err := s.GetByUsername(session.Values["username"].(string))
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}
	f := &storage.File{UserID: u.UserID}
	f.Metadata.Name = r.URL.Query().Get("filename")
	err = s.DeleteFile(f)
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File deleted")) // nolint:errcheck
}
