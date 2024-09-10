package jc

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/app"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/internal/client/util"
	jhash "github.com/JIIL07/jcloud/pkg/hash"
	"os"
	"time"
)

func AddFile(fs *app.FileService) error {
	if err := fs.F.SetFile(); err != nil {
		return fmt.Errorf("failed to prepare file info: %w", err)
	}

	fileExists, err := fs.Context.Storage.S.Exists(fs.F)
	if err != nil {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	if !fileExists {
		file := fs.F
		file.Meta.Size = len(file.Data)
		file.Status = "upload"
		err = fs.Context.Storage.S.AddFile(file)
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

	err = fs.Context.Storage.S.AddFile(file)
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

	data := util.ReadFull(f)

	meta.HashSum = jhash.Hash(string(data))

	file := &models.File{
		Meta:       meta,
		Status:     "upload",
		Data:       data,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	err = fs.Context.Storage.S.AddFile(file)
	if err != nil {
		return fmt.Errorf("failed to add file from path: %w", err)
	}
	return nil
}

// DeleteFile removes a file from the database based on its metadata.
func DeleteFile(fs *app.FileService) error {
	fs.F.Meta.Split()
	return fs.Context.Storage.S.DeleteFile(fs.F)
}

// DeleteAllFiles removes all files from the database.
func DeleteAllFiles(fs *app.FileService) error {
	return fs.Context.Storage.S.DeleteAllFiles()
}

// ListFiles retrieves a list of files from the specified table.
func ListFiles(fs *app.FileService) ([]models.File, error) {
	var files []models.File
	err := fs.Context.Storage.S.GetAllFiles(&files)
	return files, err
}

// DataInFile retrieves the file data from the database and sets it in the File struct.
func DataInFile(fs *app.FileService) error {
	fs.F.Meta.Split()

	rows, err := fs.Context.Storage.S.DB.Query(
		`SELECT data FROM local WHERE filename = ? AND extension = ?`,
		fs.F.Meta.Name,
		fs.F.Meta.Extension,
	)
	if err != nil {
		return fmt.Errorf("failed to query file data: %w", err)
	}
	defer rows.Close() // nolint:errcheck

	return util.WriteData(rows, fs.F)
}

// SearchFile searches for a file in the database and prints its metadata if found.
func SearchFile(fs *app.FileService) error {
	err := fs.Context.Storage.S.DB.Get(fs.F, `SELECT * FROM local WHERE filename = ? AND extension = ?`,
		fs.F.Meta.Name,
		fs.F.Meta.Extension)
	if err != nil {
		return err
	}

	fmt.Printf("Found: %v\n", *fs.F)
	return nil
}
