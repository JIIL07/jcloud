package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/JIIL07/cloudFiles-manager/internal/config"
)

type Storage struct {
	DB *sql.DB
}

type AboutUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Protocol string `json:"protocol"`
	Admin    bool   `json:"admin"`
}

func InitDatabase(config config.DBConfig) (*Storage, error) {
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
	var user = &AboutUser{}
	_, err := s.DB.Exec(`INSERT INTO users 
		(username, email, password, hashprotocol, admin) VALUES (?, ?, ?, ?, ?)`,
		user.Username, user.Email, user.Password, user.Protocol, user.Admin,
	)
	if err != nil {
		return fmt.Errorf("failed to save new user: %v", err)
	}
	return nil
}

func (s *Storage) GetAllUsers() ([]AboutUser, error) {
	var users []AboutUser
	rows, err := s.DB.Query(`SELECT username, email, password, hashprotocol, admin FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user = AboutUser{}
		if err := rows.Scan(&user.Username, &user.Email, &user.Password, &user.Protocol, &user.Admin); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Storage) Admin() error {
	adminEnv := os.Getenv("ADMIN_USER")
	if adminEnv == "" {
		return fmt.Errorf("no such env variable ADMIN_USER")
	}

	var admin AboutUser

	err := json.Unmarshal([]byte(adminEnv), &admin)
	if err != nil {
		return fmt.Errorf("error unmarshalling ADMIN_USER: %v", err)
	}

	_, err = s.DB.Exec(`INSERT INTO users 
		(username, email, password, hashprotocol, admin) VALUES (?, ?, ?, ?, ?)`,
		admin.Username, admin.Email, admin.Password, admin.Protocol, admin.Admin,
	)
	if err != nil {
		return fmt.Errorf("failed to save admin: %v", err)
	}

	return nil
}

func (s *Storage) Query(query string) (*sql.Rows, error) {
	return s.DB.Query(query)
}

func ParseRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		columnData := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))
		for i := range columnData {
			columnPointers[i] = &columnData[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = columnData[i]
		}
		results = append(results, row)
	}

	return results, rows.Err()
}
