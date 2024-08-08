package storage

import (
	"github.com/JIIL07/cloudFiles-manager/internal/lib/parsers"
)

func (s *Storage) GetAllFiles(u *UserData) ([]map[string]interface{}, error) {
	id, err := s.GetUserID(u)
	rows, err := s.DB.Query(`SELECT id, filename, extension, filesize, data FROM files WHERE user_id = ? ORDER BY "id" ASC`, id)
	if err != nil {
		return nil, err
	}

	return parsers.ParseRows(rows)
}

func (s *Storage) GetFile(f *File) ([]map[string]interface{}, error) {
	rows, err := s.DB.Query(`SELECT id, filename, extension, filesize, data FROM files WHERE user_id = ? AND id = ?`, f.UserID, f.Filename)
	if err != nil {
		return nil, err
	}
	return parsers.ParseRows(rows)
}
