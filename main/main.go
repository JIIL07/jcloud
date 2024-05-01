package main

import (
	"project/file"
)

func main() {
	db, _ := file.Open("sql\\files.db")
	defer db.Close()
	file.Show(db)
}
