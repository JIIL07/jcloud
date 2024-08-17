package storage

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/config"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	DB *sqlx.DB
}

func InitDatabase(c *config.Config) (SQLite, error) {
	db, err := sqlx.Open(c.Database.DriverName, c.Database.DataSourceName)
	if err != nil {
		return SQLite{DB: nil}, err
	}

	return SQLite{DB: db}, nil
}

func (s *SQLite) CreateTable() error {

	_, err := s.DB.Exec(`CREATE TABLE IF NOT EXISTS local 
	("id" INTEGER PRIMARY KEY AUTOINCREMENT, 
		"filename"  TEXT NOT NULL, 
		"extension" TEXT NOT NULL, 
		"filesize"  INTEGER NOT NULL, 
		"status" 	TEXT NOT NULL DEFAULT 'upload', 
		"data"		BLOB)`)

	return err
}

func (s *SQLite) Close() error {
	return s.DB.Close()
}

func (s *SQLite) GetAllFiles(f *[]models.File) error {
	return s.DB.Select(&f, "SELECT * FROM local")
}

func (s *SQLite) AddFile(f *models.File) error {
	_, err := s.DB.Exec(`INSERT INTO local (filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?)`,
		f.Metadata.Filename,
		f.Metadata.Extension,
		f.Metadata.Filesize,
		f.Status,
		f.Data)
	if err != nil {
		return fmt.Errorf("failed to insert file: %w", err)
	}
	return nil
}

func (s *SQLite) Exists(f *models.File) (bool, error) {
	var exists bool
	err := s.DB.Get(&exists,
		`SELECT EXISTS(SELECT 1 FROM local WHERE filename = ? AND extension = ?)`,
		f.Metadata.Filename,
		f.Metadata.Extension)

	return exists, err
}
