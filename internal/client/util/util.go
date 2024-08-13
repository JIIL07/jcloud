package util

import (
	"database/sql"
	"fmt"
)

func Find(db *sql.DB, fln, ext string) (*sql.Rows, error) {
	rows, err := db.Query(`SELECT * FROM files WHERE filename = ? AND extension = ?`, fln, ext)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		rows.Close()
		return nil, fmt.Errorf("no rows found")
	}
	return rows, nil
}

func Exists(db *sql.DB, fln, ext string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM files WHERE filename = ? AND extension = ?)`
	err := db.QueryRow(query, fln, ext).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return false, fmt.Errorf("such file already exists")
	}
	return false, nil
}

func ScanRows(rows *sql.Rows) ([]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	var results []interface{}
	for rows.Next() {
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		result := make(map[string]interface{})
		for i, colName := range columns {
			result[colName] = values[i]
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func ScanRow(rows *sql.Rows) (map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	if err = rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for i, colName := range columns {
		result[colName] = values[i]
	}
	return result, nil
}
