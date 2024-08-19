package main

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/anchor"
	"github.com/JIIL07/jcloud/internal/client/models"
)

func main() {
	files := []models.File{
		{
			ID: 1,
			Metadata: models.FileMetadata{
				Filename:  "example1.txt",
				Extension: ".txt",
				Filesize:  1234,
			},
			Status: "active",
			Data:   []byte("Hello, World!"),
		},
		{
			ID: 2,
			Metadata: models.FileMetadata{
				Filename:  "example2.jpg",
				Extension: ".jpg",
				Filesize:  5678,
			},
			Status: "active",
			Data:   []byte("Image Data"),
		},
	}

	anchorMessage := "Initial Anchor"
	a, err := anchor.NewAnchor(files, anchorMessage)
	if err != nil {
		fmt.Printf("Error during Anchor: %v\n", err)
		return
	}

	fmt.Printf("Anchor successful: %v\n", a)
}
