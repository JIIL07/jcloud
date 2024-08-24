package jc

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/app"
	"github.com/JIIL07/jcloud/internal/client/config"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/internal/client/util"
)

// AddFile inserts the file metadata and data into the database if it does not already exist.
func AddFile(fs *app.FileService) error {
	// Подготовка информации о файле
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
		file.Status = config.Statuses[0]
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
	files := []models.File{}
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
	defer rows.Close()

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
