package cloudfiles

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

func find(db *sql.DB, fln, ext string) (*sql.Rows, error) {
	rows, err := db.Query(`SELECT * FROM files WHERE filename = ? AND extension = ?`, fln, ext)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		rows.Close()
		return nil, fmt.Errorf("no rows found")
	}
	return rows, nil
}

func exists(db *sql.DB, fln, ext string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM files WHERE filename = ? AND extension = ?)`
	err := db.QueryRow(query, fln, ext).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return false, fmt.Errorf("such file already exists")
	}
	return false, nil
}

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
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open explorer: %v", err)
	}
	return nil
}

func createTempDir() (string, error) {
	tempDir, err := os.MkdirTemp("", "package_files_TEMPDIR")
	if err != nil {
		return "", fmt.Errorf("error creating temporary directory: %v", err)
	}
	return tempDir, nil
}

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
					log.Printf("Detected new file: %s", event.Name)
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
	fmt.Println(info.Filename, info.Extension)

	fileData, err := os.ReadFile(filepath.Join(tempDir, info.Fullname))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	info.Data = fileData
	info.Filesize = len(fileData)
	return info, nil
}
