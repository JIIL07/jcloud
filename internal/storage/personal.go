package storage

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/JIIL07/jcloud/pkg/bool"
)

func (s *Storage) Query(query string) (*sql.Rows, error) {
	return s.DB.Query(query)
}

func Admin(p *UserData) boolean.Wrapper {
	admin := os.Getenv("ADMIN_USER")
	if admin == "" {
		return boolean.Wrapper{Value: false}
	}

	var u *UserData
	err := json.Unmarshal([]byte(admin), &u)
	if err != nil {
		return boolean.Wrapper{Value: false}
	}

	return boolean.Wrapper{Value: p.Username == u.Username && p.Password == u.Password && p.Email == u.Email}
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

func (s *Storage) GetByUsername(username string) (UserData, error) {
	var u UserData
	err := s.DB.Get(&u, "SELECT * FROM users WHERE username = ?", username)
	return u, err
}
