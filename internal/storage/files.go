package storage

import (
	"github.com/JIIL07/jcloud/internal/lib/parsers"
)

func (s *Storage) GetAllFiles(f *File) ([]map[string]interface{}, error) {
	rows, err := s.DB.Query(`SELECT id, filename, extension, filesize, data FROM files WHERE user_id = ? ORDER BY "id" ASC`, f.UserID)
	if err != nil {
		return nil, err
	}

	return parsers.ParseRows(rows)
}

func (s *Storage) GetFile(f *File) ([]map[string]interface{}, error) {
	rows, err := s.DB.Query(`SELECT id, filename, extension, filesize, data FROM files WHERE user_id = ? AND filename = ?`, f.UserID, f.Metadata.Name)
	if err != nil {
		return nil, err
	}
	return parsers.ParseRows(rows)
}

func (s *Storage) AddFile(f *File) error {
	_, err := s.DB.Exec(`INSERT INTO files 
    	(user_id, filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?, ?)`,
		f.UserID, f.Metadata.Name, f.Metadata.Extension, f.Metadata.Size, f.Status, f.Data)
	return err
}

func (s *Storage) DeleteFile(f *File) error {
	_, err := s.DB.Exec(`DELETE FROM files WHERE user_id = ? AND filename = ?`, f.UserID, f.Metadata.Name)
	return err
}

func (s *Storage) UpdateFile(f *File) error {
	_, err := s.DB.Query(`UPDATE files SET filename = ?, extension = ?, filesize = ?, data = ? WHERE user_id = ? AND id = ?`,
		f.Metadata.Name, f.Metadata.Extension, f.Metadata.Size, f.Data, f.UserID, f.ID)
	return err
}

func (s *Storage) DeleteAllFiles(f *File) error {
	_, err := s.DB.Query(`DELETE FROM files WHERE user_id = ?`, f.UserID)
	return err
}
