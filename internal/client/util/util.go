package util

import (
	"database/sql"
	"github.com/JIIL07/jcloud/internal/client/models"
	"os"
)

func WriteData(rows *sql.Rows, info *models.File) error {
	if rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}
		info.Data = data
	}
	return nil
}

func ReadFull(f *os.File) []byte {
	stat, err := f.Stat()
	if err != nil {
		return nil
	}
	data := make([]byte, stat.Size())
	_, err = f.Read(data)
	if err != nil {
		return nil
	}
	return data

}
