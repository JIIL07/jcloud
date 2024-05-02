package file

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
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

var isValidTableName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.]*$`).MatchString
var isValidDBName = regexp.MustCompile(`^[^<>:"/\\|?*]+$`).MatchString

var colorRed = "\033[91m"
var colorReset = "\033[0m"

func Init(name string) (*sql.DB, error) {
	if !isValidDBName(name) {
		return nil, fmt.Errorf("%vinvalid DB file name: %s%v", colorRed, name, colorReset)
	}

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB, name string) error {
	if !isValidTableName(name) {
		return fmt.Errorf("%vinvalid table name: %s%v", colorRed, name, colorReset)
	}

	if _, err := db.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            filename TEXT,
            extension TEXT,
            filesize INTEGER,
            status TEXT,
            data BLOB
        )
    `, name)); err != nil {
		return fmt.Errorf("%vfailed to create table: %w%v", colorRed, err, colorReset)
	}

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS deleted (
		id INTEGER,
		filename TEXT,
		extension TEXT,
		filesize INTEGER,
		status TEXT,
		data BLOB
	)
`)
	return fmt.Errorf("%v%w%v", colorRed, err, colorReset)
}
func Add(db *sql.DB) error {
	info.GetNameExt()
	exists, err := exists(db, info.Filename, info.Extension)
	if err != nil {
		return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
	}
	if !exists {
		info.Filesize = len(info.Data)
		info.Status = Statuses[0]
		_, err = db.Exec(`INSERT INTO files (filename, extension, filesize, status, data) VALUES (?,?,?,?,?)`,
			info.Filename, info.Extension, info.Filesize, info.Status, info.Data)
		return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
	}
	return nil
}

func Delete(db *sql.DB) error {
	info.GetNameExt()
	err := db.QueryRow("SELECT * FROM files WHERE filename = ? AND extension = ?",
		info.Filename, info.Extension).Scan(&info.Id, &info.Filename,
		&info.Extension, &info.Filesize, &info.Status, &info.Data)

	if err != nil {
		return fmt.Errorf("%vquery row error: %v%v", colorRed, err, colorReset)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("%vbegin transaction error: %v%v", colorRed, err, colorReset)
	}

	_, err = tx.Exec("INSERT INTO deleted (id, filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?, ?)",
		info.Id, info.Filename, info.Extension, info.Filesize, "Deteled", info.Data)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%vinsert error: %v%v", colorRed, err, colorReset)
	}

	_, err = tx.Exec("DELETE FROM files WHERE filename = ? AND extension = ?",
		info.Filename, info.Extension)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%vdelete error: %v%v", colorRed, err, colorReset)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%vcommit transaction error: %v%v", colorRed, err, colorReset)
	}

	return nil

}

func List(db *sql.DB, tablename string) error {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tablename))
	if err != nil {
		return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
	}
	defer rows.Close()
	id := 1
	for rows.Next() {
		err := rows.Scan(new(interface{}), &info.Filename, &info.Extension, new(interface{}), new(interface{}), new(interface{}))
		if err != nil {
			return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
		}

		fmt.Printf("\033[34mid:%v \033[36m %v.%v%v\n<--------------------->\n",
			id, info.Filename, info.Extension, colorReset)
		id++
	}
	return nil
}

func DataIn(db *sql.DB) error {
	info.GetNameExt()
	rows, err := db.Query(`SELECT data FROM files WHERE filename=? AND extension=?`, info.Filename, info.Extension)
	if err != nil {
		return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&info.Data)
		if err != nil {
			return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
		}
		fmt.Printf("\033[36mData:%v \n%v\n\n", strings.TrimSpace(string(info.Data)), colorReset)
	}
	return nil
}
func Search(db *sql.DB) error {
	info.GetNameExt()
	rows, err := find(db, info.Filename, info.Extension)
	if err != nil {
		return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
	}
	err = rows.Scan(new(interface{}), &info.Filename, &info.Extension, new(interface{}), new(interface{}), new(interface{}))
	if err != nil {
		return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
	}

	fmt.Printf("\033[32m%v.%v%v\n",
		info.Filename, info.Extension, colorReset)
	return nil
}
func WriteData(db *sql.DB) error {
	info.GetNameExt()
	info.GetData()
	_, err := db.Exec(`UPDATE files SET data = ? WHERE filename = ? AND extension = ?`, info.Data, info.Filename, info.Extension)
	return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
}

func CreateFile(db *sql.DB) error {
	info.GetNameExt()
	rows, err := find(db, info.Filename, info.Extension)
	if err != nil {
		return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
	}
	if rows != nil {
		err := rows.Scan(new(interface{}), new(interface{}), new(interface{}), new(interface{}), new(interface{}), &info.Data)
		if err != nil {
			return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
		}
		file, err := os.Create(info.Fullname)
		if err != nil {
			return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
		}
		defer file.Close()
		_, err = file.Write(info.Data)
		if err != nil {
			return fmt.Errorf("%v%v%v", colorRed, err, colorReset)
		}
	}
	return nil
}
