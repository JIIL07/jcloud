package storage

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/pkg/home"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path/filepath"
)

type SQLite struct {
	DB *sqlx.DB
}

func MustInit() *SQLite {
	db, err := sqlx.Open("sqlite3", filepath.Join(home.GetHome(), ".jcloud/local_storage.db"))
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS local 
	("id" INTEGER PRIMARY KEY AUTOINCREMENT, 
		"filename"  TEXT NOT NULL, 
		"extension" TEXT NOT NULL, 
		"filesize"  INTEGER NOT NULL, 
		"status" 	TEXT NOT NULL DEFAULT 'upload', 
		"data"		BLOB)
	`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	return &SQLite{DB: db}
}

func (s *SQLite) Close() error {
	return s.DB.Close()
}

func (s *SQLite) GetAllFiles(f *[]models.File) error {
	query := `
        SELECT
            id,
            filename AS "metadata.filename",
            extension AS "metadata.extension",
            filesize AS "metadata.filesize",
            status,
            data
        FROM local
    `
	return s.DB.Select(f, query)
}
func (s *SQLite) AddFile(f *models.File) error {
	_, err := s.DB.Exec(`INSERT INTO local (filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?)`,
		f.Metadata.Name,
		f.Metadata.Extension,
		f.Metadata.Size,
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
		f.Metadata.Name,
		f.Metadata.Extension)

	return exists, err
}

func (s *SQLite) DeleteFile(f *models.File) error {
	_, err := s.DB.Exec(`DELETE FROM local WHERE filename = ? AND extension = ?`,
		f.Metadata.Name,
		f.Metadata.Extension)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *SQLite) DeleteAllFiles() error {
	_, err := s.DB.Exec(`DELETE FROM local`)
	if err != nil {
		return fmt.Errorf("failed to delete all files: %w", err)
	}
	return nil
}
