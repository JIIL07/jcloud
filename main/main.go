package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"project/file"
)

func main() {
	db, err := file.Init("files.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = file.CreateTable(db, "files")
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
				err := file.Add(db)
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}
				fmt.Println("File added successfully")

			case "list files":
				err := file.List(db, "files")
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}

			case "list deleted":
				err := file.List(db, "deleted")
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}

			case "delete":
				err := file.Delete(db)
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}
				fmt.Println("File deleted successfully")
			case "search":
				err := file.Search(db)
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}
			case "write":
				err := file.WriteData(db)
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}
				fmt.Println("File written successfully")
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
