package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/config"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

const (
	tempDirPrefix  = "jcloud-tmp"
	tempDirTimeout = 30 * time.Second
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
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	return nil
}

func CreateTempDir() (string, error) {
	dir, err := os.MkdirTemp("", tempDirPrefix)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}
	return dir, nil
}

func InitializeWatcher(tempDir string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file system watcher: %w", err)
	}

	if err := watcher.Add(tempDir); err != nil {
		return nil, fmt.Errorf("failed to add directory to watcher: %w", err)
	}

	return watcher, nil
}
func HandleFileEvents(watcher *fsnotify.Watcher, ctx context.Context) (string, error) {
	var detectedFile string
	var watcherErr error

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					detectedFile = event.Name
					return
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				watcherErr = fmt.Errorf("error watching file: %w", err)
				return
			}
		}
	}()

	wg.Wait()

	return detectedFile, watcherErr
}

func WaitForFile(tempDir string) (string, error) {
	watcher, err := InitializeWatcher(tempDir)
	if err != nil {
		return "", err
	}
	defer watcher.Close()

	ctx, cancel := context.WithTimeout(context.Background(), tempDirTimeout)
	defer cancel()

	detectedFile, watcherErr := HandleFileEvents(watcher, ctx)

	if watcherErr != nil {
		return "", watcherErr
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return "", fmt.Errorf("timeout waiting for file in directory: %s", tempDir)
	}

	return detectedFile, nil
}

// GetFileFromExplorer opens the file explorer, waits for the user to create or modify a file,
// reads the file, and returns the file's data wrapped in a models.File object.
func GetFileFromExplorer() (*models.File, error) {
	dir, err := CreateTempDir()
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}

	if err := OpenExplorer(dir); err != nil {
		return nil, fmt.Errorf("failed to open explorer: %w", err)
	}

	filePath, err := WaitForFile(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get filePath from explorer: %w", err)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read temp directory %s: %w", dir, err)
	}

	if len(files) != 1 {
		return nil, fmt.Errorf("expected exactly one file, but found %d files", len(files))
	}

	fileEntry := files[0]
	if fileEntry.IsDir() {
		return nil, fmt.Errorf("expected a file, but found a directory: %s", fileEntry.Name())
	}

	time.Sleep(time.Nanosecond * 1)

	meta := models.NewFileMetadata(fileEntry.Name())

	fileData, err := os.ReadFile(filePath)
	meta.Size = len(fileData)

	f := &models.File{
		Metadata: meta,
		Status:   config.Statuses[0],
		Data:     fileData,
	}

	return f, nil
}
