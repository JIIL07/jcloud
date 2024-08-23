package main

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/anchor"
	"github.com/JIIL07/jcloud/internal/client/delta"
	"github.com/JIIL07/jcloud/internal/client/models"
)

func main() {
	files := []models.File{
		{ID: 1, Metadata: models.FileMetadata{Name: "example1.txt", Extension: "txt", Size: 1024}, Status: "new", Data: []byte("Hello, World!")},
		{ID: 1, Metadata: models.FileMetadata{Name: "example2.txt", Extension: "txt", Size: 2048}, Status: "modified", Data: []byte("Hello, Golang!")},
	}

	previousSnapshots := make(map[int]*delta.Snapshot)

	a, err := anchor.NewAnchor(files, "Initial commit", previousSnapshots)
	if err != nil {
		fmt.Println("Error creating anchor:", err)
		return
	}

	fmt.Printf("Created Anchor: %+v\n", a)
}
