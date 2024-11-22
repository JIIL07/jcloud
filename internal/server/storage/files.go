package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type File struct {
	ID         int `db:"id"`
	UserID     int `db:"user_id"`
	Metadata   FileMetadata
	Status     string    `db:"status"`
	Data       []byte    `db:"data"`
	CreatedAt  time.Time `db:"created_at"`
	ModifiedAt time.Time `db:"last_modified_at"`
}

type FileMetadata struct {
	Name        string `db:"filename"`
	Extension   string `db:"extension"`
	Size        int    `db:"filesize"`
	HashSum     string `db:"hash_sum"`
	Description string `db:"description,omitempty"`
}

func (s *Storage) GetAllFiles(userID int) ([]File, error) {
	query := `
		SELECT id, user_id, filename AS "metadata.filename", extension AS "metadata.extension", filesize AS "metadata.filesize", 
		hash_sum AS "metadata.hash_sum", description AS "metadata.description", data, status, created_at, last_modified_at
		FROM files 
		WHERE user_id = ? 
		ORDER BY id
	`

	var files []File
	err := s.DB.Select(&files, query, userID)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *Storage) GetFile(userID int, filename string) (*File, error) {
	query := `
		SELECT id, user_id, filename AS "metadata.filename", extension AS "metadata.extension", filesize AS "metadata.filesize", 
		hash_sum AS "metadata.hash_sum", description AS "metadata.description", data, status, created_at, last_modified_at
		FROM files
		WHERE user_id = ? AND filename = ?
		ORDER BY id
	`

	var file File
	err := s.DB.Get(&file, query, userID, filename)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("error querying file: %w", err)
	}
	return &file, nil
}

func (s *Storage) AddFile(f *File) error {
	query := `
		INSERT INTO files (user_id, filename, extension, filesize, hash_sum, description, data, status, created_at, last_modified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.DB.Exec(query,
		f.UserID, f.Metadata.Name, f.Metadata.Extension, f.Metadata.Size, f.Metadata.HashSum, f.Metadata.Description,
		f.Data, f.Status, f.CreatedAt, f.ModifiedAt,
	)
	return err
}

func (s *Storage) AddFileTx(tx *sqlx.Tx, f *File) error {
	query := `
		INSERT INTO files (user_id, filename, extension, filesize, hash_sum, description, data, status, created_at, last_modified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := tx.Exec(query,
		f.UserID, f.Metadata.Name, f.Metadata.Extension, f.Metadata.Size, f.Metadata.HashSum, f.Metadata.Description,
		f.Data, f.Status, f.CreatedAt, f.ModifiedAt,
	)
	return err
}

func (s *Storage) DeleteFile(f *File) error {
	query := `DELETE FROM files WHERE user_id = ? AND filename = ? AND extension = ?`
	_, err := s.DB.Exec(query, f.UserID, f.Metadata.Name, f.Metadata.Extension)
	return err
}

func (s *Storage) DeleteAllFiles(userID int) error {
	query := `DELETE FROM files WHERE user_id = ?`
	_, err := s.DB.Exec(query, userID)
	return err
}

func (s *Storage) RenameFile(f *File) error {
	query := `
			UPDATE files 
			SET filename = ?
			WHERE user_id = ? AND id = ?
	`
	_, err := s.DB.Exec(query, f.Metadata.Name)
	return err
}

func (s *Storage) GetImageFiles(userID int) ([]File, error) {
	query := `
		SELECT id, user_id, filename AS "metadata.filename", extension AS "metadata.extension", filesize AS "metadata.filesize", 
		hash_sum AS "metadata.hash_sum", description AS "metadata.description", data, status, created_at, last_modified_at
		FROM files
		WHERE user_id = ? AND extension IN ('jpg', 'jpeg', 'png')
	`

	var files []File
	err := s.DB.Select(&files, query, userID)
	if err != nil {
		return nil, err
	}

	return files, nil
}
