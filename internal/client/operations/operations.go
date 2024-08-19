package operations

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/config"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/internal/client/storage"
	"github.com/JIIL07/jcloud/internal/client/util"
)

type FileContext struct {
	Info    *models.File
	Storage *storage.SQLite
}

// AddFile inserts the file metadata and data into the database if it does not already exist.
func (fctx *FileContext) AddFile() error {
	if err := fctx.Info.SetFile(); err != nil {
		return fmt.Errorf("failed to prepare info: %w", err)
	}

	fileExists, err := fctx.Storage.Exists(fctx.Info)
	if err != nil {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	if !fileExists {
		fctx.Info.Metadata.Filesize = len(fctx.Info.Data)
		fctx.Info.Status = config.Statuses[0]
		err = fctx.Storage.AddFile(fctx.Info)
		if err != nil {
			return fmt.Errorf("failed to add file: %w", err)
		}
	}
	return nil
}

func (fctx *FileContext) AddFileFromExplorer() error {
	f, err := util.GetFileFromExplorer()
	if err != nil {
		return fmt.Errorf("failed to get file from explorer: %w", err)
	}

	fctx.Info = f
	err = fctx.Storage.AddFile(fctx.Info)
	return err
}

// DeleteFile removes a file from the database based on its metadata.
func (fctx *FileContext) DeleteFile() error {
	fctx.Info.Metadata.Split()

	_, err := fctx.Storage.DB.Exec(`DELETE FROM files WHERE filename = ? AND extension = ?`,
		fctx.Info.Metadata.Filename, fctx.Info.Metadata.Extension)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// ListFiles retrieves a list of files from the specified table.
func (fctx *FileContext) ListFiles() ([]models.File, error) {
	files := &[]models.File{}
	err := fctx.Storage.GetAllFiles(files)
	return *files, err
}

// DataInFile retrieves the file data from the database and sets it in the Info struct.
func (fctx *FileContext) DataInFile() error {
	fctx.Info.Metadata.Split()
	rows, err := fctx.Storage.DB.Query(`SELECT data FROM local WHERE filename = ? AND extension = ?`,
		fctx.Info.Metadata.Filename,
		fctx.Info.Metadata.Extension)
	if err != nil {
		return fmt.Errorf("failed to query file data: %w", err)
	}
	defer rows.Close()

	// Assuming WriteData processes the rows to set the file data in File
	return util.WriteData(rows, fctx.Info)
}

// SearchFile searches for a file in the database and prints its metadata if found.
func (fctx *FileContext) SearchFile() error {
	err := fctx.Storage.DB.Get(fctx.Info, `SELECT * FROM local WHERE filename = ? AND extension = ?`,
		fctx.Info.Metadata.Filename,
		fctx.Info.Metadata.Extension)
	if err != nil {
		return err
	}

	fmt.Printf("Found: %v\n", *fctx.Info)
	return nil
}
