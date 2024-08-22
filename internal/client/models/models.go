package models

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

type FileMetadata struct {
	Name      string `db:"filename"`
	Extension string `db:"extension"`
	Size      int    `db:"filesize"`
}

func (metadata *FileMetadata) Split() {
	dotIndex := strings.LastIndex(metadata.Name, ".")
	if dotIndex != -1 {
		metadata.Extension = metadata.Name[dotIndex+1:]
		metadata.Name = metadata.Name[:dotIndex]
	} else {
		metadata.Extension = ""
	}
}

func NewFileMetadata(fullname string) FileMetadata {
	metadata := FileMetadata{Name: fullname}
	metadata.Split()
	return metadata
}

type File struct {
	ID       int `db:"id"`
	Metadata FileMetadata
	Status   string `db:"status"`
	Data     []byte `db:"data"`
}

func (i *File) SetFile() error {
	if i == nil {
		return errors.New("info struct is nil")
	}

	m, err := ReadNameFromStdin()
	if err != nil {
		return fmt.Errorf("error reading name: %v", err)
	}
	i.Metadata = NewFileMetadata(m)

	data, err := ReadDataFromStdin()
	if err != nil {
		return fmt.Errorf("error reading data: %v", err)
	}
	i.Data = data

	return nil
}

func ReadNameFromStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the full file name: ")

	fullname, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading full name from stdin: %v", err)
	}

	fullname = strings.TrimSpace(fullname)
	if len(fullname) == 0 {
		return "", errors.New("full name cannot be empty")
	}

	return fullname, nil
}

func ReadDataFromStdin() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Data to read (Press Ctrl+D to end reading)")

	data, err := reader.ReadBytes('\x04')
	if err != nil {
		return nil, fmt.Errorf("error reading data from stdin: %v", err)
	}

	data = bytes.TrimSuffix(data, []byte{'\x04'})

	return data, nil
}
