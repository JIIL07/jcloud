package main

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/models/v2/builder"
	"log"

	"github.com/JIIL07/jcloud/internal/client/models"
)

func main() {
	info := &models.File{}
	err := info.SetFile()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File: %+v\n", info)
}

func Builder(id int) {
	b := &builder.InfoBuilder{}
	b.WithID(id)

	fullname, err := models.ReadNameFromStdin()
	if err != nil {
		log.Fatal(err)
	}

	metadata := models.NewFileMetadata(fullname)
	b.WithMetadata(metadata)

	data, err := models.ReadDataFromStdin()
	if err != nil {
		log.Fatal(err)
	}

	b.WithData(data)
	b.WithStatus("active")

	info := b.Build()

	fmt.Printf("File: %+v\n", info)
}
