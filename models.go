package cloudfiles

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Info struct {
	Id        int
	Filename  string
	Extension string
	Filesize  int
	Status    string
	Data      []byte
	Fullname  string
}

func (info *Info) PrepareInfo() error {
	if err := info.GetNameExt(); err != nil {
		return err
	}

	return info.GetData()
}

func (info *Info) Split() {
	parts := strings.Split(info.Fullname, ".")

	if len(parts) > 1 {
		info.Filename = strings.Join(parts[:len(parts)-1], ".")
		info.Extension = parts[len(parts)-1]
	} else {
		fmt.Println("\033[33mWARNING: empty extension\033[0m")
		info.Filename = info.Fullname
		info.Extension = ""
	}
}

func (info *Info) GetNameExt() error {
	var err error

	info.Fullname, err = readFullName()
	if err != nil {
		return fmt.Errorf("error reading file name: %v", err)
	}

	info.Split()

	return nil
}

func (info *Info) GetData() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Data to read \033[32m(Press Ctrl+D to end reading)\033[0m: ")

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
