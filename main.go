package main

import (
	"log"
)

func main() {
	dbPath := "sql\\files.db"
	db, err := sqlOpen(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
