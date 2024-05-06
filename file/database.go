package file

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Database describes general methods of working with a database.
type Database interface {
	Init(name string) (*sql.DB, error)
	CreateTable(db *sql.DB, name string) error
}

// SQLiteDB implements the Database interface for SQLite.
type SQLiteDB struct{}

// Init initializes a connection to the SQLite database.
func (s *SQLiteDB) Init(name string) (*sql.DB, error) {
	if !isValidDBName(name) {
		return nil, fmt.Errorf("invalid DB file name: %s", name)
	}
	return sql.Open("sqlite3", name)
}

// CreateTable creates a table with the given name if it does not exist.
func (s *SQLiteDB) CreateTable(db *sql.DB, name string) error {
	if !isValidTableName(name) {
		return fmt.Errorf("invalid table name: %s", name)
	}
	createStmt := `CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT,
		extension TEXT,
		filesize INTEGER,
		status TEXT,
		data BLOB
	)`
	_, err := db.Exec(fmt.Sprintf(createStmt, name))
	return err
}
