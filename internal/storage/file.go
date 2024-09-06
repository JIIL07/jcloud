package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func (s *Storage) GetAllFiles(userID int) ([]File, error) {
	var files []File
	query := `SELECT id, filename AS "metadata.filename", extension AS "metadata.extension", filesize AS "metadata.filesize", data, user_id, status
		FROM files 
		WHERE user_id = ? 
		ORDER BY id`
	err := s.DB.Select(&files, query, userID)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *Storage) GetFile(userID int, filename string) (*File, error) {
	var file File
	query := `SELECT id, filename AS "metadata.filename", extension AS "metadata.extension", filesize AS "metadata.filesize", data, user_id, status
		FROM files
		WHERE user_id = ? AND filename = ?`
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
	query := `INSERT INTO files (user_id, filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.DB.Exec(query, f.UserID, f.Metadata.Name, f.Metadata.Extension, f.Metadata.Size, f.Status, f.Data)
	return err
}

func (s *Storage) AddFileTx(tx *sqlx.Tx, f *File) error {
	_, err := tx.Exec(`INSERT INTO files (user_id, filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?, ?)`,
		f.UserID, f.Metadata.Name, f.Metadata.Extension, f.Metadata.Size, f.Status, f.Data)
	return err
}

func (s *Storage) DeleteFile(f *File) error {
	query := `DELETE FROM files WHERE user_id = ? AND filename = ?`
	_, err := s.DB.Exec(query, f.UserID, f.Metadata.Name)
	return err
}

func (s *Storage) DeleteAllFiles(userID int) error {
	query := `DELETE FROM files WHERE user_id = ?`
	_, err := s.DB.Exec(query, userID)
	return err
}

func (s *Storage) UpdateFile(f *File) error {
	query := `UPDATE files SET filename = ?, extension = ?, filesize = ?, data = ? WHERE user_id = ? AND id = ?`
	_, err := s.DB.Exec(query, f.Metadata.Name, f.Metadata.Extension, f.Metadata.Size, f.Data, f.UserID, f.ID)
	return err
}
