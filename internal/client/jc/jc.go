package jc

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/config"
	"github.com/JIIL07/jcloud/internal/client/lib/home"
	jlog "github.com/JIIL07/jcloud/internal/client/lib/logger"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/internal/client/storage"
	"github.com/JIIL07/jcloud/internal/client/util"
	"log/slog"
)

type Client struct {
	cfg    *config.Config
	common service
}

type service struct {
	client  *Client
	file    *models.File
	storage *storage.SQLite
	paths   *home.Paths
	logger  *slog.Logger
}

var c *Client

func Init() *Client {
	c = &Client{}
	c.cfg = config.MustLoad()
	c.common.client = c
	c.common.file = &models.File{}
	c.common.storage = storage.MustInit(c.cfg)
	c.common.paths = home.SetPaths()
	c.common.logger = jlog.NewLogger(c.common.paths.Jlog)

	return c
}

// AddFile inserts the file metadata and data into the database if it does not already exist.
func (s *service) AddFile() error {
	if err := s.file.SetFile(); err != nil {
		return fmt.Errorf("failed to prepare info: %w", err)
	}
	fileExists, err := s.client.ctx.Storage.Exists(s.client.ctx.File)
	if err != nil {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	if !fileExists {
		fctx.File.Metadata.Size = len(fctx.File.Data)
		fctx.File.Status = config.Statuses[0]
		err = fctx.Storage.AddFile(fctx.File)
		if err != nil {
			return fmt.Errorf("failed to add file: %w", err)
		}
	}
	return nil
}

func (fctx *Context) AddFileFromExplorer() error {
	f, err := util.GetFileFromExplorer()
	if err != nil {
		return fmt.Errorf("failed to get file from explorer: %w", err)
	}

	fctx.File = f
	err = fctx.Storage.AddFile(fctx.File)
	return err
}

// DeleteFile removes a file from the database based on its metadata.
func (fctx *Context) DeleteFile() error {
	fctx.File.Metadata.Split()
	return fctx.Storage.DeleteFile(fctx.File)
}

// DeleteAllFiles removes all files from the database.
func (fctx *Context) DeleteAllFiles() error {
	return fctx.Storage.DeleteAllFiles()
}

// ListFiles retrieves a list of files from the specified table.
func (fctx *Context) ListFiles() ([]models.File, error) {
	files := &[]models.File{}
	err := fctx.Storage.GetAllFiles(files)
	return *files, err
}

// DataInFile retrieves the file data from the database and sets it in the File struct.
func (fctx *Context) DataInFile() error {
	fctx.File.Metadata.Split()
	rows, err := fctx.Storage.DB.Query(`SELECT data FROM local WHERE filename = ? AND extension = ?`,
		fctx.File.Metadata.Name,
		fctx.File.Metadata.Extension)
	if err != nil {
		return fmt.Errorf("failed to query file data: %w", err)
	}
	defer rows.Close()

	// Assuming WriteData processes the rows to set the file data in File
	return util.WriteData(rows, fctx.File)
}

// SearchFile searches for a file in the database and prints its metadata if found.
func (fctx *Context) SearchFile() error {
	err := fctx.Storage.DB.Get(fctx.File, `SELECT * FROM local WHERE filename = ? AND extension = ?`,
		fctx.File.Metadata.Name,
		fctx.File.Metadata.Extension)
	if err != nil {
		return err
	}

	fmt.Printf("Found: %v\n", *fctx.File)
	return nil
}
