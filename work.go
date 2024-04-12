package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func splitName(name string) (string, string) {
	split := strings.Split(name, ".")
	return split[0], split[1]
}
func getName() error {
	fmt.Print("Filename: ")
	_, err := fmt.Scan(&info.Filename)
	return err
}

func getID(db *sql.DB) error {
	var lastID int
	err := db.QueryRow("SELECT MAX(id) FROM files").Scan(&lastID)
	info.Id = lastID + 1
	return err
}
func dirCreate() (string, error) {
	dir := "./sql"
	err := os.Mkdir(dir, 0755)
	dbPath := filepath.Join(dir, "files.db")
	return dbPath, err
}
func Search(db *sql.DB) (*sql.Rows, error) {
	query, filename, extension := Form()
	row, err := db.Query(query, filename, extension)
	return row, err
}
func Form() (string, string, string) {
	fmt.Print("File to search: ")
	search.fullNotation, _ = reader.ReadString('\n')
	search.fullNotation = strings.TrimSpace(search.fullNotation)
	search.name, search.ext = splitName(search.fullNotation)
	return `SELECT * FROM files WHERE filename = ? AND extension = ? `, search.name, "." + search.ext
}
func Write(*sql.Rows) {

}
