package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/JIIL07/cloudFiles-manager/internal/config"
)

type Storage struct {
	DB *sql.DB
}

type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Protocol string `json:"protocol"`
	Admin    int    `json:"admin"`
}

func InitDatabase(config *config.Config) (*Storage, error) {
	if config.Env == "prod" || config.Env == "debug" {
		config.Database.DataSourceName = os.Getenv("DATABASE_PATH")
	}
	db, err := sql.Open(config.Database.DriverName, config.Database.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"username" TEXT NOT NULL UNIQUE,
		"email" TEXT NOT NULL,
		"password" TEXT NOT NULL,
		"hashprotocol" TEXT,
		"admin" INTEGER DEFAULT 1
	);`)

	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) CloseDatabase() error {
	return s.DB.Close()
}

func (s *Storage) SaveNewUser(user *UserData) error {
	_, err := s.DB.Exec(`INSERT INTO users 
		(username, email, password, hashprotocol, admin) VALUES (?, ?, ?, ?, ?)`,
		user.Username, user.Email, user.Password, user.Protocol, user.Admin,
	)
	if err != nil {
		return fmt.Errorf("failed to save new user: %v", err)
	}
	return nil
}

func (s *Storage) GetAllUsers() ([]UserData, error) {
	var users []UserData
	var user = UserData{}
	rows, err := s.DB.Query(`SELECT username, email, password, hashprotocol, admin FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&user.Username, &user.Email, &user.Password, &user.Protocol, &user.Admin); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Storage) Query(query string) (*sql.Rows, error) {
	return s.DB.Query(query)
}
