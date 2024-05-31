package main

import (
	"log"

	cloud "github.com/JIIL07/cloudFiles-manager/client"
	server "github.com/JIIL07/cloudFiles-manager/server"
)

var sqlite *cloud.SQLiteDB

func main() {
	db, err := sqlite.PrepareLocalDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	servctx := server.NewServerContext(db)
	if err := servctx.Start(); err != nil {
		log.Fatal(err)
	}

}
