package file

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
)

// find retrieves a set of rows from the 'files' table matching the specified filename and extension.
func find(db *sql.DB, fln, ext string) (*sql.Rows, error) {
	query := `SELECT * FROM files WHERE filename = ? AND extension = ?`
	rows, err := db.Query(query, fln, ext)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	if !rows.Next() {
		rows.Close()
		return nil, fmt.Errorf("no files found")
	}
	return rows, nil
}

// exists checks whether a file with the given filename and extension already exists in the database.
func exists(db *sql.DB, fln, ext string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM files WHERE filename = ? AND extension = ?)`
	var exists bool
	err := db.QueryRow(query, fln, ext).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %v", err)
	}
	return exists, nil
}

// openExplorer opens the file explorer based on the operating system.
func openExplorer(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}

// createTempDir creates a temporary directory for file operations.
func createTempDir() (string, error) {
	tempDir, err := os.MkdirTemp("", "package_files_TEMPDIR")
	if err != nil {
		return "", fmt.Errorf("error creating temporary directory: %v", err)
	}
	return tempDir, nil
}

// waitFile waits for a new file to be created in the specified directory using file system notifications.
func waitFile(tempDir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					done <- true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Error watching file: %v", err)
			}
		}
	}()

	if err := watcher.Add(tempDir); err != nil {
		return fmt.Errorf("failed to add directory to watcher: %v", err)
	}

	<-done
	return nil
}

// processFile processes the newly created file in the specified directory.
func processFile(tempDir string) (*Info, error) {
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read temp directory: %v", err)
	}

	if len(files) != 1 {
		return nil, fmt.Errorf("expected exactly one file, got %d", len(files))
	}

	fileEntry := files[0]
	if fileEntry.IsDir() {
		return nil, fmt.Errorf("a file was expected, not a directory")
	}

	info := &Info{Fullname: fileEntry.Name()}
	info.Split()

	fileData, err := os.ReadFile(filepath.Join(tempDir, info.Fullname))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	info.Data = fileData
	info.Filesize = len(fileData)
	return info, nil
}
