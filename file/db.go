package file

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Info struct {
	Id        int
	Filename  string
	Extension string
	Filesize  int
	Status    string
	Data      []byte

	Fullname string
}
type TempInfo struct {
	fullNotation string
	name         string
	ext          string
}

var Statuses = []string{
	"Created",
	"Has data in",
	"Renamed",
	"Deleted",
}

var info Info
var temp TempInfo
var err error

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	func() error {
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS deleted (
				id INTEGER,
				filename TEXT,
				extension TEXT,
				filesize INTEGER,
	        	status TEXT,
				data BLOB
			)
		`)
		if err != nil {
			return err
		}
		return err
	}()
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
	err := db.QueryRow("SELECT * FROM files WHERE id = ?", id).Scan(&info.Id,
		&info.Filename, &info.Extension, &info.Filesize, &info.Status, &info.Data)
	if err != nil {
		return fmt.Errorf("query row error: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction error: %v", err)
	}

	_, err = tx.Exec("INSERT INTO deleted (id, filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?, ?)",
		info.Id, info.Filename, info.Extension, info.Filesize, "Deteled", info.Data)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("insert error: %v", err)
	}

	_, err = tx.Exec("DELETE FROM files WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("delete error: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction error: %v", err)
	}

	return nil

}
func Show(db *sql.DB, tablename string) error {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tablename))
	if err != nil {
		return err
	}
	defer rows.Close()
	id := 1
	for rows.Next() {
		err := rows.Scan(new(interface{}), &info.Filename, &info.Extension, new(interface{}), new(interface{}), &info.Data)
		if err != nil {
			return err
		}

		fmt.Printf("\033[34mid:%v \033[36m %v.%v\033[0m\nData:\n%v\n<--------------------->\n",
			id, info.Filename, info.Extension, strings.TrimSpace(string(info.Data)))
		id++
	}
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
