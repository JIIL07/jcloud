package cloudfiles

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

var db *sql.DB
var s *SQLiteDB
var ctx *FileContext

func InitDB() *sql.DB {
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
	return db
}

func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	items, err := ctx.List("files")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
