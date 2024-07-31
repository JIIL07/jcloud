package storage

import "database/sql"

func ParseRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		columnData := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))
		for i := range columnData {
			columnPointers[i] = &columnData[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = columnData[i]
		}
		results = append(results, row)
	}

	return results, rows.Err()
}
