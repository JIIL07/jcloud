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
		"admin" TEXT
	);`)

	return &Storage{DB: db}, err
}
