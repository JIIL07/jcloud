package main

import (
	file "Dev/project/fileAPI"
)

func main() {
	db, _ := file.Open("sql\\files.db")
	defer db.Close()
	file.Show(db)
}
