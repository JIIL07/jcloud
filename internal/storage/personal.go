package storage

import (
	"database/sql"
	"encoding/json"
	"os"
)

func (s *Storage) Query(query string) (*sql.Rows, error) {
	return s.DB.Query(query)
}

func Admin(p *UserData) bool {
	admin := os.Getenv("ADMIN_USER")
	if admin == "" {
		return false
	}

	var u *UserData
	err := json.Unmarshal([]byte(admin), &u)
	if err != nil {
		return false
	}

	return p.Username == u.Username && p.Password == u.Password && p.Email == u.Email
}

func (s *Storage) CheckExistence(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`
	err := s.DB.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Storage) GetUserID(u *UserData) (int, error) {
	var id int
	err := s.DB.Get(
		&id,
		`SELECT id FROM users WHERE username = ?`,
		u.Username)
	return id, err
}
func (s *Storage) GetByUsername(username string) (*UserData, error) {
	var u UserData
	err := s.DB.Get(&u, "SELECT * FROM users WHERE username = ?", username)
	return &u, err
}
