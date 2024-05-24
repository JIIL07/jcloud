package cloudfiles

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	Init() (*sql.DB, error)
	CreateTable(db *sql.DB, name string) error
}

type SQLiteDB struct{}

func (s *SQLiteDB) Init() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (s *SQLiteDB) CreateTable(db *sql.DB, name string) error {
	if !isValidTableName(name) {
		return fmt.Errorf("invalid table name: %s", name)
	}
	_, err := db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
	(id INTEGER PRIMARY KEY AUTOINCREMENT, 
		filename TEXT, 
		extension TEXT, 
		filesize INTEGER, 
		status TEXT, 
		data BLOB)`, name))
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS deleted 
	(id INTEGER PRIMARY KEY AUTOINCREMENT, 
		filename TEXT, 
		extension TEXT, 
		filesize INTEGER, 
		status TEXT, 
		data BLOB)`)
	return err
}
