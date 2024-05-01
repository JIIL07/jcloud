package file

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
)

func Open(path string) (*sql.DB, error) {
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
func Add(db *sql.DB) error {
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
func Delete(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM files WHERE id = ?", id)
	return err
}
func Show(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM files")
	if err != nil {
		return err
	}
	defer rows.Close()

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"id", "filename", "extension", "size", "status"})
	for rows.Next() {
		err := rows.Scan(&info.Id, &info.Filename, &info.Extension, &info.Filesize, &info.Status, new(interface{}))
		if err != nil {
			return err
		}
		values := []interface{}{info.Id, info.Filename, info.Extension, info.Filesize, info.Status}
		data := make([]string, len(values))
		for i, v := range values {
			data[i] = fmt.Sprintf("%v", v)
		}
		table.Append(data)
	}
	table.Render()
	return nil
}
func Search(db *sql.DB) (*sql.Rows, error) {
	GetName("search")
	rows, err := Find(db, temp.name, temp.ext)
	return rows, err
}
func CreateFile(db *sql.DB) error {
	GetName("file create")
	rows, err := Find(db, temp.name, temp.ext)
	if err != nil {
		return err
	}
	if rows != nil {
		err := rows.Scan(new(interface{}), new(interface{}), new(interface{}), new(interface{}), new(interface{}), &info.Data)
		if err != nil {
			return err
		}
		file, err := os.Create(temp.fullNotation)
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
