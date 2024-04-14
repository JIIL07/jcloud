package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func sqlOpen(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT ,
			filename TEXT,
			extension TEXT,
			filesize INTEGER,
            status TEXT,
			data BLOB
		)
	`)
	return db, err
}
func sqlAdd(db *sql.DB) error {
	GetName("add")
	err = GetData()
	if err != nil {
		return err
	}
	info.Filesize = len(info.Data)
	info.Status = Statuses[0]
	_, err = db.Exec(`INSERT INTO files
	(filename, extension, filesize, status, data)
	VALUES
	(?,?,?,?,?)`, info.Filename, info.Extension, info.Filesize, info.Status, info.Data)
	return err
}
func sqlDelete(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM files WHERE id = ?", id)
	return err
}
func sqlShow(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM files")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&info.Id, &info.Filename, &info.Extension, &info.Filesize, &info.Status, &info.Data)
		if err != nil {
			return err
		}
		fmt.Printf("%d %s %s %d %s %s\n", info.Id, info.Filename, info.Extension, info.Filesize, info.Status, info.Data)
	}
	return nil
}
func sqlSearch(db *sql.DB) (*sql.Rows, error) {
	GetName("search")
	rows, err := Find(db, search.name, search.ext)
	return rows, err
}
func sqlCreateFile(db *sql.DB) error {
	GetName("file create")
	rows, err := Find(db, createFile.name, createFile.ext)
	if err != nil {
		return err
	}
	if rows != nil {
		err := rows.Scan(new(interface{}), new(interface{}), new(interface{}), new(interface{}), new(interface{}), &info.Data)
		if err != nil {
			return err
		}
		file, err := os.Create(createFile.fullNotation)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.Write(info.Data)
		if err != nil {
			return err
		}

	}
	return nil
}
