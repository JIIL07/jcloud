package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type item struct {
	Path     string `json:"path"`
	Selected bool   `json:"selected"`
}

func runRustApp() error {
	cmd := exec.Command("./interactive.exe")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func readSelectedFilesList(filename string) ([]item, error) {
	var files []item
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func rustInteractive() {
	err := runRustApp()
	if err != nil {
		fmt.Println("Error running Rust app:", err)
		return
	}

	updatedFiles, err := readSelectedFilesList("selected.json")
	if err != nil {
		fmt.Println("Error reading files from JSON:", err)
		return
	}

	for _, file := range updatedFiles {
		if file.Selected {
			addFile(file.Path)
		}
	}
}
