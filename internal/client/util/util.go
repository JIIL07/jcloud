package util

import (
	"database/sql"
	"github.com/JIIL07/jcloud/internal/client/models"
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
