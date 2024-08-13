package util

import (
	"fmt"
	"github.com/JIIL07/cloudFiles-manager/internal/client/models"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func OpenExplorer(path string) error {
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

func CreateTempDir() (string, error) {
	tempDir, err := os.MkdirTemp("", "package_files_TEMPDIR")
	if err != nil {
		return "", fmt.Errorf("error creating temporary directory: %v", err)
	}
	return tempDir, nil
}

func WaitFile(tempDir string) error {
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

func ProcessFile(tempDir string) (*models.Info, error) {
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

	meta := models.NewFileMetadata(fileEntry.Name())

	fileData, err := os.ReadFile(filepath.Join(tempDir, fileEntry.Name()))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	meta.Filesize = len(fileData)
	info := &models.Info{Metadata: meta}
	info.Data = fileData

	return info, nil
}
