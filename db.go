package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func sqlOpen(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY ,
			filename TEXT,
			extension TEXT,
			filesize INTEGER,
            status TEXT,
			text TEXT
		)
	`)
	return db, err
}
func sqlAdd(db *sql.DB, file []string) error {
	info.Filename = file[0]
	info.Extension = "." + file[1]
	info.Filesize = len(info.Text)
	info.Status = Statuses[0]
	_, err := db.Exec(`INSERT INTO files
	(id, filename, extension, filesize, status, text)
	VALUES
	(?,?,?,?,?,?)`, info.Id, info.Filename, info.Extension, info.Filesize, info.Status, info.Text)
	return err
}
func sqlDelete(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM files WHERE id = ?", id)
	return err
}
func sqlShow(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM files LIMIT -1 OFFSET 1")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var info Info
		err := rows.Scan(&info.Id, &info.Filename, &info.Extension, &info.Filesize, &info.Status, &info.Text)
		if err != nil {
			return err
		}
		fmt.Printf("%d %s %s %d %s %s\n", info.Id, info.Filename, info.Extension, info.Filesize, info.Status, info.Text)
	}
	return nil
}
