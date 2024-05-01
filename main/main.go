package main

import (
	"fmt"
	"log"
	"project/file"
)

func main() {
	db, err := file.Open("files.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	file.Show(db, "files")
	fmt.Println("DELETED FILES")
	file.Show(db, "deleted")
}
