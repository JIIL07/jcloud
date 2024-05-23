package cloudfiles

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

var db *sql.DB
var s *SQLiteDB

func InitDB() {
	var err error
	db, err = s.Init()
	if err != nil {
		log.Fatal(err)
	}
	err = s.CreateTable(db, "files")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name FROM files")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		items = append(items, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
