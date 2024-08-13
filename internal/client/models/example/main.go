package main

import (
	"fmt"
	"log"

	"github.com/JIIL07/cloudFiles-manager/internal/client/models"
)

func main() {
	info := &models.Info{}
	fullname, err := models.ReadNameFromStdin()
	if err != nil {
		log.Fatal(err)
	}

	metadata := models.NewFileMetadata(fullname)
	info.Metadata = metadata
	err = info.SetData()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Info: %+v\n", info)
}

func Builder(id int) {
	builder := &models.InfoBuilder{}
	builder.WithID(id)

	fullname, err := models.ReadNameFromStdin()
	if err != nil {
		log.Fatal(err)
	}

	metadata := models.NewFileMetadata(fullname)
	builder.WithMetadata(metadata)

	data, err := models.ReadDataFromStdin()
	if err != nil {
		log.Fatal(err)
	}

	builder.WithData(data)
	builder.WithStatus("active")

	info := builder.Build()

	fmt.Printf("Info: %+v\n", info)
}
