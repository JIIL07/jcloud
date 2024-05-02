package file

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strings"
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

type Splitter interface {
	Split()
}

type Getter interface {
	GetNameExt()
	GetData()
}

func (info *Info) Split() {
	parts := strings.Split(info.Fullname, ".")
	if len(parts) > 1 {
		info.Filename = parts[0]
		info.Extension = parts[len(parts)-1]
	} else {
		fmt.Println("\033[93mWARNING: empty extension\033[0m")
		info.Filename = info.Fullname
		info.Extension = ""
	}
}
func (info *Info) GetNameExt() {
	info.Fullname = readFullName()
	info.Split()
}

func (info *Info) GetData() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Data to read \033[32m(Press Ctrl+D to end reading)\033[0m: ")
	data, err := reader.ReadBytes('\x04')
	if err != nil {
		fmt.Println("Error reading data:", err)
	}
	info.Data = bytes.TrimSuffix(data, []byte{'\x04'})
}

func readFullName() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the full name of the file: ")
	fullname, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading fullname:", err)
		return ""
	}
	return strings.TrimSpace(fullname)
}

func find(db *sql.DB, fln, ext string) (*sql.Rows, error) {
	rows, err := db.Query(`SELECT * FROM files WHERE filename = ? AND extension = ? `, fln, ext)
	if !rows.Next() {
		fmt.Println("No rows found")
		return nil, err
	}
	return rows, err
}
func exists(db *sql.DB, fln, ext string) (bool, error) {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT * FROM files WHERE filename = ? AND extension = ?)`, fln, ext).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, fmt.Errorf("such file already exists")
	}
	return exists, nil
}
