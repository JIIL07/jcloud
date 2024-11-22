package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/JIIL07/jcloud/pkg/bool"
)

type User struct {
	UserID   int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Protocol string `db:"hashprotocol" json:"hashprotocol"`
	Admin    int    `db:"admin" json:"admin"`
}

func Admin(p *User) boolean.Wrapper {
	admin := os.Getenv("ADMIN_USER")
	if admin == "" {
		return boolean.Wrapper{Value: false}
	}

	var u *User
	err := json.Unmarshal([]byte(admin), &u)
	if err != nil {
		return boolean.Wrapper{Value: false}
	}

	return boolean.Wrapper{Value: p.Username == u.Username && p.Password == u.Password && p.Email == u.Email}
}

func (s *Storage) CheckExistence(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`
	err := s.DB.Get(&exists, query, username)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Storage) GetByUsername(username string) (User, error) {
	var u User
	query := `SELECT * FROM users WHERE username = ?`
	err := s.DB.Get(&u, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, fmt.Errorf("user not found")
		}
		return u, fmt.Errorf("error querying user: %w", err)
	}
	return u, nil
}

func (s *Storage) SaveNewUser(u *User) error {
	query := `
		INSERT INTO users (username, email, password, hashprotocol, admin) 
		VALUES (:username, :email, :password, :protocol, :admin)
	`
	_, err := s.DB.NamedExec(query, u)
	if err != nil {
		return fmt.Errorf("failed to save new user: %v", err)
	}
	return nil
}

func (s *Storage) GetAllUsers() ([]User, error) {
	var users []User
	query := `SELECT * FROM users`
	err := s.DB.Select(&users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %v", err)
	}
	return users, nil
}

func (s *Storage) DeleteUser(username string) error {
	query := `DELETE FROM users WHERE username = ?`
	_, err := s.DB.Exec(query, username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

func (s *Storage) UpdateUser(username string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	var setClauses []string
	args := make([]interface{}, 0, len(updates)+1)

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}

	query := fmt.Sprintf(`UPDATE users SET %s WHERE username = ?`, strings.Join(setClauses, ", "))
	args = append(args, username)

	_, err := s.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}
