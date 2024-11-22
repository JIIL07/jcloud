package home

import (
	"fmt"
	"os"
	"path/filepath"
)

type Paths struct {
	Home       string
	JcloudFile *os.File
	JlogFile   *os.File
	Jcookie    *os.File
	AnchorFile *os.File
	Profile    string
}

func GetHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return ""
	}
	return home
}

func createDir(baseDir, subDir string) string {
	dirPath := filepath.Join(baseDir, subDir)
	err := os.MkdirAll(dirPath, 0750)
	if err != nil {
		fmt.Printf("Error creating directory %s: %v\n", dirPath, err)
		return ""
	}
	return dirPath
}

func createFile(dir, filename string) *os.File {
	filePath := filepath.Join(dir, filename)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		fmt.Printf("Error creating/opening file %s: %v\n", filePath, err)
		return nil
	}
	return file
}

func SetPaths() *Paths {
	homeDir := GetHome()
	if homeDir == "" {
		return nil
	}

	jcloudDir := createDir(homeDir, ".jcloud")
	jlogDir := createDir(jcloudDir, ".jlog")
	anchorDir := createDir(jcloudDir, ".anchor")

	jcloudFile := createFile(jcloudDir, ".jcloud")
	logFile := createFile(jlogDir, "jlog.log")
	jcookieFile := createFile(jcloudDir, ".jcookie")
	anchorLog := createFile(anchorDir, "anchor.log")

	return &Paths{
		Home:       homeDir,
		JcloudFile: jcloudFile,
		JlogFile:   logFile,
		Jcookie:    jcookieFile,
		AnchorFile: anchorLog,
	}
}

func (p *Paths) Close() {
	files := []*os.File{p.JcloudFile, p.JlogFile, p.Jcookie, p.AnchorFile}
	for _, file := range files {
		if file != nil {
			if err := file.Close(); err != nil {
				fmt.Printf("Error closing file: %v\n", err)
			}
		}
	}
}
