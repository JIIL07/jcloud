package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	file "github.com/JIIL07/cloudFiles-manager"
)

var dbHandler file.SQLiteDB

func main() {
	db, err := dbHandler.Init("files.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := file.FileContext{
		DB:   db,
		Info: &file.Info{},
	}

	err = dbHandler.CreateTable(db, "files")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(os.Stdin)

outerLoop:
	for {
		fmt.Println("Enter command:")
	innerLoop:
		for scanner.Scan() {
			command := scanner.Text()
			switch command {
			case "add":
				err := ctx.Add()
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}
				fmt.Println("File added successfully")

			case "addfile":
				err := ctx.AddFile()
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}
				fmt.Println("File added successfully")

			case "list files":
				err := ctx.List("files")
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}

			case "list deleted":
				err := ctx.List("deleted")
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}

			case "delete":
				err := ctx.Delete()
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}
				fmt.Println("File deleted successfully")
			case "search":
				err := ctx.Search()
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}

			case "write":
				err := ctx.WriteData()
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}
				fmt.Println("File written successfully")

			case "datain":
				err := ctx.DataIn()
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}

			case "exit":
				break outerLoop

			default:
				fmt.Println("Unknown command")
			}
			fmt.Println("Enter command:")
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from input:", err)
		}
	}

}
