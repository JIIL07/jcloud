package server

import (
	"encoding/json"
	"fmt"
	"github.com/JIIL07/jcloud/internal/storage"
	jctx "github.com/JIIL07/jcloud/pkg/ctx"
	"net/http"
	"strings"
)

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	u, ok := jctx.FromContext[*storage.UserData](r.Context(), "user")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}
	s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}

	files, err := s.GetAllFiles(u.UserID)
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
	u, ok := jctx.FromContext[*storage.UserData](r.Context(), "user")
	if !ok {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}

	var files []storage.File
	err := json.NewDecoder(r.Body).Decode(&files)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tx, err := s.DB.Beginx()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				http.Error(w, "Failed to rollback transaction: "+err.Error(), http.StatusInternalServerError)
			}
			http.Error(w, "Failed to add files: "+err.Error(), http.StatusInternalServerError)
		}
	}()

	for _, file := range files {
		file.UserID = u.UserID
		err = s.AddFileTx(tx, &file)
		if err != nil {
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Files added successfully")) // nolint:errcheck
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	u, ok := jctx.FromContext[*storage.UserData](r.Context(), "user")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}

	s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}

	f := &storage.File{UserID: u.UserID}
	f.Metadata.Name = r.URL.Query().Get("filename")
	err := s.DeleteFile(f)
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File deleted")) // nolint:errcheck
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	u, ok := jctx.FromContext[*storage.UserData](r.Context(), "user")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}
	s, ok := jctx.FromContext[*storage.Storage](r.Context(), "storage")
	if !ok {
		http.Error(w, "Storage not found", http.StatusInternalServerError)
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	file, err := s.GetFile(u.UserID, strings.Split(filename, ".")[0])
	if err != nil {
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(file.Data)))
	w.Header().Set("Content-Disposition", "attachment; filename="+file.Metadata.Name+"."+file.Metadata.Extension)
	w.Header().Set("Content-Type", "application/octet-stream")

	w.Write(file.Data) // nolint:errcheck
}
