package storage

import (
	"database/sql"
	"fmt"
	"github.com/JIIL07/jcloud/internal/server/config"
	"github.com/jmoiron/sqlx"
	"os"
)

type Storage struct {
	DB *sqlx.DB
}

func InitDatabase(config *config.Config) (*Storage, error) {
	if config.Env == "prod" || config.Env == "debug" {
		config.Database.DataSourceName = os.Getenv("DATABASE_PATH")
	}
	db, err := sqlx.Open(config.Database.DriverName, config.Database.DataSourceName)
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

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS files (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"user_id" INTEGER NOT NULL,
		"filename" TEXT NOT NULL, 
		"extension" TEXT NOT NULL, 
		"filesize" INTEGER NOT NULL, 
		"status" TEXT, 
		"data" BLOB,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"last_modified_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"hash_sum" TEXT NOT NULL,
		"description" TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id),
		UNIQUE(user_id, filename)
	);`)

	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) CloseDatabase() error {
	return s.DB.Close()
}

func (s *Storage) Query(command string) (*sql.Rows, error) {
	return s.DB.Query(command)
}
