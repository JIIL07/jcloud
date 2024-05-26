package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	cloudfiles "github.com/JIIL07/cloudFiles-manager"
)

var (
	ctx    *cloudfiles.FileContext
	server *cloudfiles.ServerContext
	mu     sync.Mutex
)

func main() {

	sqliteDB := &cloudfiles.SQLiteDB{}
	serverDB, err := sqliteDB.PrepareLocalDB()
	if err != nil {
		log.Fatal(err)
	}

	server = cloudfiles.NewServerContext(serverDB)

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}

	localDB, err := sqliteDB.PrepareLocalDB()
	if err != nil {
		log.Fatal(err)
	}

	ctx = &cloudfiles.FileContext{
		DB:   localDB,
		Info: &cloudfiles.Info{},
	}

	scanner := bufio.NewScanner(os.Stdin)

outerLoop:
	for {
		fmt.Println("Enter command:")
	innerLoop:
		for scanner.Scan() {
			command := scanner.Text()
			mu.Lock()
			switch command {
			case "add":
				err := ctx.Add()
				if err != nil {
					fmt.Println(err)
					mu.Unlock()
					break innerLoop
				}
				fmt.Println("File added successfully")

			case "addfile":
				err := ctx.AddFile()
				if err != nil {
					fmt.Println(err)
					mu.Unlock()
					break innerLoop
				}
				fmt.Println("File added successfully")

			case "list files":
				items, err := ctx.List("files")
				if err != nil {
					fmt.Println(err)
					mu.Unlock()
					break innerLoop
				}
				fmt.Println(items)

			case "list deleted":
				items, err := ctx.List("deleted")
				if err != nil {
					fmt.Println(err)
					mu.Unlock()
					break innerLoop
				}
				fmt.Println(items)

			case "delete":
				err := ctx.Delete()
				if err != nil {
					fmt.Println(err)
					mu.Unlock()
					break innerLoop
				}
				fmt.Println("File deleted successfully")

			case "search":
				err := ctx.Search()
				if err != nil {
					fmt.Println(err)
					mu.Unlock()
					break innerLoop
				}

			case "serversearch":
				err := server.Ctx.Search()
				if err != nil {
					fmt.Println(err)
					mu.Unlock()
					break innerLoop
				}

			case "write":
				err := ctx.WriteData()
				if err != nil {
					fmt.Println(err)
					mu.Unlock()
					break innerLoop
				}
				fmt.Println("File written successfully")

			case "datain":
				err := ctx.DataIn()
				if err != nil {
					fmt.Println(err)
					mu.Unlock()
					break innerLoop
				}

			case "push":
				err := cloudfiles.CopyDB(localDB, serverDB)
				if err != nil {
					fmt.Println("Error pushing data to server:", err)
					mu.Unlock()
					break innerLoop
				}
				fmt.Println("Data pushed to server successfully")

			case "exit":
				mu.Unlock()
				break outerLoop

			default:
				fmt.Println("Unknown command")
			}
			mu.Unlock()
			fmt.Println("Enter command:")
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from input:", err)
		}
	}
}
