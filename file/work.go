package file

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
)

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
