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
func (ctx *FileContext) AddFile() error {
	if err := ctx.Info.SetFile(); err != nil {
		return fmt.Errorf("failed to prepare info: %w", err)
	}

	fileExists, err := ctx.Storage.Exists(ctx.Info)
	if err != nil {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	if !fileExists {
		ctx.Info.Metadata.Filesize = len(ctx.Info.Data)
		ctx.Info.Status = config.Statuses[0]
		err = ctx.Storage.AddFile(ctx.Info)
		if err != nil {
			return fmt.Errorf("failed to add file: %w", err)
		}
	}
	return nil
}

func (ctx *FileContext) AddFileFromExplorer() error {
	f, err := util.GetFileFromExplorer()
	if err != nil {
		return fmt.Errorf("failed to get file from explorer: %w", err)
	}

	ctx.Info = f
	err = ctx.Storage.AddFile(ctx.Info)
	return err
}

// DeleteFile removes a file from the database based on its metadata.
func (ctx *FileContext) DeleteFile() error {
	ctx.Info.Metadata.Split()

	_, err := ctx.Storage.DB.Exec(`DELETE FROM files WHERE filename = ? AND extension = ?`,
		ctx.Info.Metadata.Filename, ctx.Info.Metadata.Extension)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// ListFiles retrieves a list of files from the specified table.
func (ctx *FileContext) ListFiles() ([]models.File, error) {
	files := &[]models.File{}
	err := ctx.Storage.GetAllFiles(files)
	return *files, err
}

// DataInFile retrieves the file data from the database and sets it in the Info struct.
func (ctx *FileContext) DataInFile() error {
	ctx.Info.Metadata.Split()
	rows, err := ctx.Storage.DB.Query(`SELECT data FROM local WHERE filename = ? AND extension = ?`,
		ctx.Info.Metadata.Filename,
		ctx.Info.Metadata.Extension)
	if err != nil {
		return fmt.Errorf("failed to query file data: %w", err)
	}
	defer rows.Close()

	// Assuming WriteData processes the rows to set the file data in File
	return util.WriteData(rows, ctx.Info)
}

// SearchFile searches for a file in the database and prints its metadata if found.
func (ctx *FileContext) SearchFile() error {
	err := ctx.Storage.DB.Get(ctx.Info, `SELECT * FROM local WHERE filename = ? AND extension = ?`,
		ctx.Info.Metadata.Filename,
		ctx.Info.Metadata.Extension)
	if err != nil {
		return err
	}

	fmt.Printf("Found: %v\n", *ctx.Info)
	return nil
}
