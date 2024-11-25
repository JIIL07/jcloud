package jc

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/app"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/internal/client/util"
	jhash "github.com/JIIL07/jcloud/pkg/hash"
	"log"
	"os"
	"time"
)

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
	defer f.Close() // nolint:errcheck

	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stat: %w", err)
	}

	meta := models.NewFileMetadata(stat.Name())

	data := util.ReadFull(f)
	var cBuf bytes.Buffer
	gzipWriter := gzip.NewWriter(&cBuf)
	_, err = gzipWriter.Write(data)
	if err != nil {
		log.Fatal("Error compressing data:", err)
	}
	gzipWriter.Close() // nolint:errcheck

	meta.Size = len(cBuf.Bytes())
	meta.HashSum = jhash.Hash(data)

	file := &models.File{
		Meta:       meta,
		Status:     "upload",
		Data:       cBuf.Bytes(),
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
