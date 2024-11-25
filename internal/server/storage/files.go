package storage

import (
	"database/sql"
	"errors"
	"fmt"
	jhash "github.com/JIIL07/jcloud/pkg/hash"
	"github.com/jmoiron/sqlx"
	"time"
)

type File struct {
	ID         int `db:"id" json:"id"`
	UserID     int `db:"user_id" json:"user_id"`
	Metadata   FileMetadata
	Status     string    `db:"status" json:"status"`
	Data       []byte    `db:"data" json:"data"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	ModifiedAt time.Time `db:"last_modified_at" json:"modified_at"`
}

type FileMetadata struct {
	Name        string `db:"filename" json:"filename"`
	Extension   string `db:"extension" json:"extension"`
	Size        int    `db:"filesize" json:"filesize"`
	HashSum     string `db:"hash_sum" json:"hash_sum"`
	Description string `db:"description,omitempty" json:"description,omitempty"`
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

func (s *Storage) UpdateFile(f *File, newData []byte) error {
	params := map[string]interface{}{
		"userID":     f.UserID,
		"filename":   f.Metadata.Name,
		"status":     "modified",
		"modifiedAt": time.Now(),
		"data":       newData,
		"size":       len(newData),
		"hashSum":    jhash.Hash(newData),
	}

	query := `
		UPDATE files
		SET 
			status = :status,
			last_modified_at = :modifiedAt,
			data = :data, 
			filesize = :size, 
			hash_sum = :hashSum
		WHERE user_id = :userID AND filename = :filename
	`

	_, err := s.DB.NamedExec(query, params)
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	return nil
}

func (s *Storage) UpdateFileMetadata(userID int, req struct {
	Filename    string `json:"filename"`
	Extension   string `json:"extension"`
	Description string `json:"description"`
	OldName     string `json:"oldname"`
}) error {

	params := map[string]interface{}{
		"userID":      userID,
		"oldName":     req.OldName,
		"filename":    req.Filename,
		"extension":   req.Extension,
		"description": req.Description,
		"status":      "metadata modified",
		"modifiedAt":  time.Now(),
	}

	query := `
		UPDATE files
		SET 
		    filename = :filename,
			extension = :extension,
			description = :description,
			status = :status,
			last_modified_at = :modifiedAt
		WHERE filename = :oldName AND user_id = :userID
	`

	_, err := s.DB.NamedExec(query, params)
	if err != nil {
		return fmt.Errorf("failed to update metadata: %w", err)
	}
	return nil
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
