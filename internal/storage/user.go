package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, fmt.Errorf("user not found")
		}
		return u, fmt.Errorf("error querying user: %w", err)
	}
	return u, nil
}

func (s *Storage) SaveNewUser(u *UserData) error {
	_, err := s.DB.Exec(`INSERT INTO users 
		(username, email, password, hashprotocol, admin) VALUES (?, ?, ?, ?, ?)`,
		u.Username, u.Email, u.Password, u.Protocol, u.Admin,
	)
	if err != nil {
		return fmt.Errorf("failed to save new user: %v", err)
	}
	return nil
}

func (s *Storage) GetAllUsers() ([]UserData, error) {
	var users []UserData
	var u = &UserData{}
	rows, err := s.DB.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&u.Username, &u.Email, &u.Password, &u.Protocol, &u.Admin); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, *u)
	}
	return users, nil
}

func (s *Storage) DeleteUser(username string) error {
	_, err := s.DB.Exec(`DELETE FROM users WHERE username = ?`, username)
	if err != nil {
		return err
	}
	return nil
}
