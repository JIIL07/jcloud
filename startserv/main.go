package main

import (
	"log"

	cloudfiles "github.com/JIIL07/cloudFiles-manager"
)

var sqlite *cloudfiles.SQLiteDB

func main() {
	db, err := sqlite.PrepareLocalDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := cloudfiles.NewServerContext(db)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
