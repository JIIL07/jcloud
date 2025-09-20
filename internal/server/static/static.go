package static

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type Files map[string][]byte

func (f Files) BinaryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	data, exists := f[filename]
	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(data)
}

func LoadStatic(path string) (*Files, error) {
	staticFiles := Files{}

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			filename := filepath.Base(filePath)
			staticFiles[filename] = data
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load static files: %w", err)
	}

	return &staticFiles, nil
}
