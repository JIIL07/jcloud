package storage

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	DB *sqlx.DB
}

func InitDatabase(c *config.Config) (SQLite, error) {
	db, err := sqlx.Open(c.Database.DriverName, c.Database.DataSourceName)
	if err != nil {
		return SQLite{DB: nil}, err
	}

	return SQLite{DB: db}, nil
}

func (s *SQLite) CreateTable(name string) error {
	if !config.IsValidTableName(name) {
		return fmt.Errorf("invalid table name: %s", name)
	}
	_, err := s.DB.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
	(id INTEGER PRIMARY KEY AUTOINCREMENT, 
		filename TEXT, 
		extension TEXT, 
		filesize INTEGER, 
		status TEXT, 
		data BLOB)`, name))
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(`CREATE TABLE IF NOT EXISTS deletedfiles 
	(id INTEGER PRIMARY KEY AUTOINCREMENT, 
		filename TEXT, 
		extension TEXT, 
		filesize INTEGER, 
		status TEXT, 
		data BLOB)
	`)

	return err
}
