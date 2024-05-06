package file

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// Info holds metadata about a file.
type Info struct {
	Id        int
	Filename  string
	Extension string
	Filesize  int
	Status    string
	Data      []byte
	Fullname  string
}

// PrepareInfo prepares the file information for database operations.
func (info *Info) PrepareInfo() error {
	if err := info.GetNameExt(); err != nil {
		return err
	}
	return info.GetData()
}

// Split extracts the filename and extension from Fullname.
func (info *Info) Split() {
	parts := strings.Split(info.Fullname, ".")
	if len(parts) > 1 {
		info.Filename = strings.Join(parts[:len(parts)-1], ".")
		info.Extension = parts[len(parts)-1]
	} else {
		fmt.Println("WARNING: empty extension")
		info.Filename = info.Fullname
		info.Extension = ""
	}
}

// GetNameExt retrieves the full filename and extension.
func (info *Info) GetNameExt() error {
	fullname, err := readFullName()
	if err != nil {
		return fmt.Errorf("error reading full name: %v", err)
	}
	info.Fullname = fullname
	info.Split()
	return nil
}

// GetData reads file data from standard input for simulation purposes.
func (info *Info) GetData() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Data to read (Press Ctrl+D to end reading): ")
	data, err := reader.ReadBytes('\x04')
	if err != nil {
		return fmt.Errorf("error reading data: %v", err)
	}
	info.Data = bytes.TrimSuffix(data, []byte{'\x04'})
	return nil
}

func readFullName() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the full file name: ")
	fullname, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading full name: %v", err)
	}
	return strings.TrimSpace(fullname), nil
}
