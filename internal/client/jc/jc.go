package jc

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/app"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/internal/client/util"
	jlog "github.com/JIIL07/jcloud/pkg/log"
	"os"
	"path/filepath"
)

func AddFile(fs *app.FileService) error {
	if err := fs.F.SetFile(); err != nil {
		return fmt.Errorf("failed to prepare file info: %w", err)
	}

	fileExists, err := fs.Context.StorageService.S.Exists(fs.F)
	if err != nil {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	if !fileExists {
		file := fs.F
		file.Metadata.Size = len(file.Data)
		file.Status = "upload"
		err = fs.Context.StorageService.S.AddFile(file)
		if err != nil {
			return fmt.Errorf("failed to add file: %w", err)
		}
	}
	return nil
}

func AddFileFromExplorer(fs *app.FileService) error {
	file, err := util.GetFileFromExplorer()
	if err != nil {
		return fmt.Errorf("failed to get file from explorer: %w", err)
	}

	err = fs.Context.StorageService.S.AddFile(file)
	if err != nil {
		return fmt.Errorf("failed to add file from explorer: %w", err)
	}
	return nil
}

func AddFileFromPath(fs *app.FileService, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file from path: %w", err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stat: %w", err)
	}

	meta := models.NewFileMetadata(stat.Name())
	meta.Size = int(stat.Size())
	file := &models.File{
		Metadata: meta,
		Status:   "upload",
		Data:     util.ReadFull(f),
	}

	err = fs.Context.StorageService.S.AddFile(file)
	if err != nil {
		return fmt.Errorf("failed to add file from path: %w", err)
	}
	return nil
}

func AddFilesFromDir(fs *app.FileService, dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fs.Context.LoggerService.L.Error("Error accessing path", jlog.Err(err), "path", path)
			return fmt.Errorf("failed to access path %s: %w", path, err)
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			fs.Context.LoggerService.L.Error("Failed to open file", jlog.Err(err), "file", path)
			return fmt.Errorf("failed to open file at path %s: %w", path, err)
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %w", path, err)
		}

		meta := models.NewFileMetadata(stat.Name())
		meta.Size = int(stat.Size())

		newFile := &models.File{
			Metadata: meta,
			Status:   "upload",
			Data:     util.ReadFull(file),
		}

		err = fs.Context.StorageService.S.AddFile(newFile)
		if err != nil {
			fs.Context.LoggerService.L.Error("Failed to add file to storage", jlog.Err(err), "file", path)
			return fmt.Errorf("failed to add file %s to storage: %w", path, err)
		}

		fs.Context.LoggerService.L.Info("Successfully added file", "file", path)
		return nil
	})

	if err != nil {
		fs.Context.LoggerService.L.Error("Failed to walk through directory", jlog.Err(err), "directory", dirPath)
		return fmt.Errorf("failed to add files from directory %s: %w", dirPath, err)
	}

	return nil
}

// DeleteFile removes a file from the database based on its metadata.
func DeleteFile(fs *app.FileService) error {
	fs.F.Metadata.Split()
	return fs.Context.StorageService.S.DeleteFile(fs.F)
}

// DeleteAllFiles removes all files from the database.
func DeleteAllFiles(fs *app.FileService) error {
	return fs.Context.StorageService.S.DeleteAllFiles()
}

// ListFiles retrieves a list of files from the specified table.
func ListFiles(fs *app.FileService) ([]models.File, error) {
	var files []models.File
	err := fs.Context.StorageService.S.GetAllFiles(&files)
	return files, err
}

// DataInFile retrieves the file data from the database and sets it in the File struct.
func DataInFile(fs *app.FileService) error {
	fs.F.Metadata.Split()

	rows, err := fs.Context.StorageService.S.DB.Query(
		`SELECT data FROM local WHERE filename = ? AND extension = ?`,
		fs.F.Metadata.Name,
		fs.F.Metadata.Extension,
	)
	if err != nil {
		return fmt.Errorf("failed to query file data: %w", err)
	}
	defer rows.Close() // nolint:errcheck

	return util.WriteData(rows, fs.F)
}

// SearchFile searches for a file in the database and prints its metadata if found.
func SearchFile(fs *app.FileService) error {
	err := fs.Context.StorageService.S.DB.Get(fs.F, `SELECT * FROM local WHERE filename = ? AND extension = ?`,
		fs.F.Metadata.Name,
		fs.F.Metadata.Extension)
	if err != nil {
		return err
	}

	fmt.Printf("Found: %v\n", *fs.F)
	return nil
}
