package main

import (
	"fmt"
	"log"
)

func main() {
	dbPath, _ := dirCreate()
	db, err := sqlOpen(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := Search(db)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&info.Id, &info.Filename, &info.Extension, &info.Filesize, &info.Status, &info.Text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d %s%s %d %s %s\n", info.Id, info.Filename, info.Extension, info.Filesize, info.Status, info.Text)
	}
}
