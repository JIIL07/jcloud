package storage

import (
	"database/sql"
	"fmt"
)

type DBConfig struct {
	DriverName     string `yaml:"driverName" env-required:"true"`
	DataSourceName string `yaml:"dataSourceName" env-required:"true"`
}

type Storage struct {
	DB *sql.DB
}

type User struct {
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Protocol string `json:"Protocol"`
	Admin    bool   `json:"Admin"`
}

func InitDatabase(config DBConfig) (*Storage, error) {
	db, err := sql.Open(config.DriverName, config.DataSourceName)
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
		"admin" BOOLEAN DEFAULT FALSE
	);`)

	return &Storage{DB: db}, err
}

func (s *Storage) CloseDatabase() error {
	return s.DB.Close()
}

func (s *Storage) SaveNewUser() error {
	var user = &User{}
	_, err := s.DB.Exec(`INSERT INTO users 
		(username, email, password, hashprotocol, admin) VALUES (?, ?, ?, ?, ?)`,
		user.Username, user.Email, user.Password, user.Protocol, user.Admin,
	)
	if err != nil {
		return fmt.Errorf("failed to save new user: %v", err)
	}
	return nil
}

func (s *Storage) GetAllUsers() ([]User, error) {
	var users []User
	rows, err := s.DB.Query(`SELECT username, email, password, hashprotocol, admin FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user = User{}
		if err := rows.Scan(&user.Username, &user.Email, &user.Password, &user.Protocol, &user.Admin); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}
