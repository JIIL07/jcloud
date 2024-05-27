package main

import (
	cmd "github.com/JIIL07/cloudFiles-manager/cmd"
	//server "github.com/JIIL07/cloudFiles-manager/server"
)

func main() {
	cmd.Execute()

	// server := server.NewServerContext(db)

	// if err := server.Start(); err != nil {
	// 	log.Fatal(err)
	// }
}
