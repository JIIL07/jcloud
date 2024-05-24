package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	file "github.com/JIIL07/cloudFiles-manager"
)

var db *sql.DB

func main() {
	db = file.InitDB()

	var ctx = file.FileContext{
		DB:   db,
		Info: &file.Info{},
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
				items, err := ctx.List("files")
				if err != nil {
					fmt.Println(err)
					break innerLoop
				}
				fmt.Println(items)

			case "list deleted":
				_, err := ctx.List("deleted")
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

	openServer(&ctx)
}

func openServer(ctx *file.FileContext) {

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		items, err := ctx.List("files")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}
