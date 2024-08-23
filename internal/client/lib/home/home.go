package home

import (
	"os"
	"path/filepath"
)

type Paths struct {
	Home    string
	Jcloud  *os.File
	Jlog    *os.File
	Profile string
}

func GetHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

func CreateJcloudDir(home string) string {
	jcloudDir := filepath.Join(home, ".jcloud")
	err := os.MkdirAll(jcloudDir, os.ModePerm)
	if err != nil {
		return ""
	}
	return jcloudDir
}

func CreateJlogDir(jcloudDir string) string {
	jlogDir := filepath.Join(jcloudDir, ".jlog")
	err := os.MkdirAll(jlogDir, os.ModePerm)
	if err != nil {
		return ""
	}
	return jlogDir
}

func CreateJcloudFile(jcloudDir string) *os.File {
	file, err := os.OpenFile(filepath.Join(jcloudDir, ".jcloud"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}
	return file
}

func CreateLogFile(jlogDir string) *os.File {
	file, err := os.OpenFile(filepath.Join(jlogDir, "jcloud.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}
	return file
}

func SetPaths() *Paths {
	homeDir := GetHome()
	jcloudDir := CreateJcloudDir(homeDir)
	jlogDir := CreateJlogDir(jcloudDir)

	jcloudFile := CreateJcloudFile(jcloudDir)
	logFile := CreateLogFile(jlogDir)

	return &Paths{
		Home:   homeDir,
		Jcloud: jcloudFile,
		Jlog:   logFile,
	}
}
