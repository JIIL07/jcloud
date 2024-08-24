package main

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/anchor"
	"github.com/JIIL07/jcloud/internal/client/delta"
	"github.com/JIIL07/jcloud/internal/client/lib/home"
	"github.com/JIIL07/jcloud/internal/client/models"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	files := []models.File{}

	for i := 1; i <= 50; i++ {
		file := models.File{
			ID: i % 2,
			Metadata: models.FileMetadata{
				Name:      "example" + strconv.Itoa(i),
				Extension: "txt",
				Size:      1024,
			},
			Status: "new",
			Data:   []byte("Hello, Golang! " + strconv.Itoa(i)),
		}
		files = append(files, file)
	}

	previousSnapshots := make(map[int]*delta.Snapshot)

	a, err := anchor.NewAnchor(files, "Initial commit", previousSnapshots)
	if err != nil {
		fmt.Println("Error creating anchor:", err)
		return
	}

	err = os.WriteFile(filepath.Join(home.GetHome(), ".jcloud", ".anchor", "anchor.log"), []byte(a.Log), 0644)
	if err != nil {
		fmt.Println("Error writing anchor log:", err)
		return
	}
}
