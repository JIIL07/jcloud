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

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS local (
		"id"                INTEGER PRIMARY KEY AUTOINCREMENT, 
		"filename"          TEXT NOT NULL, 
		"extension"         TEXT NOT NULL, 
		"filesize"          INTEGER NOT NULL, 
		"status"            TEXT NOT NULL DEFAULT 'upload',
		"data"              BLOB,
		"created_at"        TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
		"last_modified_at"  TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
		"hash_sum"          TEXT NOT NULL,
		"description"       TEXT
	)
`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		log.Fatal(err)
	}

	return &SQLite{DB: db}
}

func (s *SQLite) Close() error {
	return s.DB.Close()
}

func (s *SQLite) GetFile(f *models.File) error {
	return s.DB.Get(f, `
		SELECT id, filename AS "m.filename", extension AS "m.extension", filesize AS "m.filesize",
		       status, data, created_at, last_modified_at, hash_sum AS "m.hash_sum", description AS "m.description"
		FROM local
		WHERE filename = ? AND extension = ?`,
		f.Meta.Name, f.Meta.Extension)
}

func (s *SQLite) GetAllFiles(f *[]models.File) error {
	return s.DB.Select(f, `
		SELECT id, filename AS "m.filename", extension AS "m.extension", filesize AS "m.filesize",
		       status, data, created_at, last_modified_at, hash_sum AS "m.hash_sum", description AS "m.description"
		FROM local`)
}

func (s *SQLite) AddFile(f *models.File) error {
	_, err := s.DB.Exec(`
		INSERT INTO local (filename, extension, filesize, status, data, created_at, last_modified_at, hash_sum, description) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		f.Meta.Name,
		f.Meta.Extension,
		f.Meta.Size,
		f.Status,
		f.Data,
		f.CreatedAt,
		f.ModifiedAt,
		f.Meta.HashSum,
		f.Meta.Description)
	if err != nil {
		return fmt.Errorf("failed to insert file: %w", err)
	}
	return nil
}

func (s *SQLite) Exists(f *models.File) (bool, error) {
	var exists bool
	err := s.DB.Get(&exists, `
		SELECT EXISTS(SELECT 1 FROM local WHERE filename = ? AND extension = ?)`,
		f.Meta.Name, f.Meta.Extension)

	return exists, err
}

func (s *SQLite) DeleteFile(f *models.File) error {
	_, err := s.DB.Exec(`DELETE FROM local WHERE filename = ? AND extension = ?`,
		f.Meta.Name, f.Meta.Extension)
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

func (s *SQLite) UpdateFile(f *models.File) error {
	_, err := s.DB.Exec(`
		UPDATE local
		SET status = ?, data = ?, last_modified_at = ?
		WHERE filename = ? AND extension = ?`,
		f.Status,
		f.Data,
		f.ModifiedAt,
		f.Meta.Name,
		f.Meta.Extension)
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}
	return nil
}

func (s *SQLite) UpdateFileDescription(f *models.File) error {
	_, err := s.DB.Exec(`
		UPDATE local
		SET description = ?
		WHERE filename = ? AND extension = ?`,
		f.Meta.Description,
		f.Meta.Name,
		f.Meta.Extension)
	if err != nil {
		return fmt.Errorf("failed to update file description: %w", err)
	}
	return nil
}
