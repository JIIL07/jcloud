package file

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var Statuses = []string{
	"Created",
	"Has data in",
	"Renamed",
	"Deleted",
}

var info Info
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
	info.GetNameExt()
	exists, err := exists(db, info.Filename, info.Extension)
	if err != nil {
		return err
	}
	if !exists {
		info.Filesize = len(info.Data)
		info.Status = Statuses[0]
		_, err = db.Exec(`INSERT INTO files (filename, extension, filesize, status, data) VALUES (?,?,?,?,?)`,
			info.Filename, info.Extension, info.Filesize, info.Status, info.Data)
		return err
	}
	return nil
}
func Delete(db *sql.DB) error {
	info.GetNameExt()
	err := db.QueryRow("SELECT * FROM files WHERE filename = ? AND extension = ?",
		info.Filename, info.Extension).Scan(&info.Id, &info.Filename,
		&info.Extension, &info.Filesize, &info.Status, &info.Data)

	if err != nil {
		return fmt.Errorf("\033[91mquery row error: %v\033[0m", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("\033[91mbegin transaction error: %v\033[0m", err)
	}

	_, err = tx.Exec("INSERT INTO deleted (id, filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?, ?)",
		info.Id, info.Filename, info.Extension, info.Filesize, "Deteled", info.Data)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("\033[91minsert error: %v\033[0m", err)
	}

	_, err = tx.Exec("DELETE FROM files WHERE filename = ? AND extension = ?",
		info.Filename, info.Extension)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("\033[91mdelete error: %v\033[0m", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("\033[91mcommit transaction error: %v\033[0m", err)
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
func Search(db *sql.DB) error {
	info.GetNameExt()
	rows, err := find(db, info.Filename, info.Extension)
	if err != nil {
		return err
	}
	err = rows.Scan(new(interface{}), &info.Filename, &info.Extension, new(interface{}), new(interface{}), &info.Data)
	if err != nil {
		return err
	}

	fmt.Printf("\033[32m%v.%v\033[0m\nData:\n%v\n",
		info.Filename, info.Extension, strings.TrimSpace(string(info.Data)))
	return nil
}
func WriteData(db *sql.DB) error {
	info.GetNameExt()
	info.GetData()
	_, err = db.Exec(`UPDATE files SET data = ? WHERE filename = ? AND extension = ?`, info.Data, info.Filename, info.Extension)
	return err
}

func CreateFile(db *sql.DB) error {
	info.GetNameExt()
	rows, err := find(db, info.Filename, info.Extension)
	if err != nil {
		return err
	}
	if rows != nil {
		err := rows.Scan(new(interface{}), new(interface{}), new(interface{}), new(interface{}), new(interface{}), &info.Data)
		if err != nil {
			return err
		}
		file, err := os.Create(info.Fullname)
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
