package storage

import (
	"time"
)

type FileVersion struct {
	FileID      int       `db:"file_id"`
	UserID      int       `db:"user_id"`
	Version     int       `db:"version"`
	FullVersion bool      `db:"full_version"`
	Delta       []byte    `db:"delta"`
	ChangeType  string    `db:"change_type"`
	CreatedAt   time.Time `db:"created_at"`
}

func (s *Storage) AddFileVersion(version FileVersion) error {
	_, err := s.DB.NamedExec(`
        INSERT INTO file_versions (file_id, user_id, version, full_version, delta, change_type)
        VALUES (:file_id, :user_id, :version, :full_version, :delta, :change_type)
    `, version)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(`
        UPDATE files
        SET last_version_id = (
            SELECT id FROM file_versions
            WHERE file_id = ? AND version = ?
        )
        WHERE id = ?
    `, version.FileID, version.Version, version.FileID)

	return err
}

func (s *Storage) RestoreFileToVersion(fileID int, targetVersion int) ([]byte, error) {
	var versions []FileVersion
	err := s.DB.Select(&versions, `
        SELECT * FROM file_versions
        WHERE file_id = ? AND version <= ?
        ORDER BY version ASC
    `, fileID, targetVersion)
	if err != nil {
		return nil, err
	}

	var fileContent []byte
	for _, version := range versions {
		fileContent, err = ApplyVersionToContent(fileContent, version)
		if err != nil {
			return nil, err
		}
	}

	return fileContent, nil
}

func (s *Storage) GetFileHistory(fileID int) ([]FileVersion, error) {
	var history []FileVersion
	err := s.DB.Select(&history, `
        SELECT version, change_type, created_at, user_id FROM file_versions
        WHERE file_id = ?
        ORDER BY version ASC
    `, fileID)
	return history, err
}

func (s *Storage) GetFileVersions(fileID int) ([]FileVersion, error) {
	var versions []FileVersion
	err := s.DB.Select(&versions, `
		SELECT * FROM file_versions
		WHERE file_id = ?
		ORDER BY version ASC
	`, fileID)
	return versions, err
}

func (s *Storage) GetFileVersion(fileID int, version int) (FileVersion, error) {
	var versionData FileVersion
	err := s.DB.Get(&versionData, `
		SELECT * FROM file_versions
		WHERE file_id = ? AND version = ?
	`, fileID, version)
	return versionData, err
}

func (s *Storage) GetLastFileVersion(fileID int) (FileVersion, error) {
	var versionData FileVersion
	err := s.DB.Get(&versionData, `
		SELECT * FROM file_versions
		WHERE file_id = ?
		ORDER BY version DESC
		LIMIT 1
	`, fileID)
	return versionData, err
}

func (s *Storage) DeleteFileVersion(fileID int, version int) error {
	_, err := s.DB.Exec(`
		DELETE FROM file_versions
		WHERE file_id = ? AND version = ?
	`, fileID, version)
	return err
}

func (s *Storage) DeleteFileVersions(fileID int) error {
	_, err := s.DB.Exec(`
		DELETE FROM file_versions
		WHERE file_id = ?
	`, fileID)
	return err
}
