package home

import (
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
		return ""
	}
	return home
}

func CreateJcloudDir(home string) string {
	jcloudDir := filepath.Join(home, ".jcloud")
	err := os.MkdirAll(jcloudDir, 0750)
	if err != nil {
		return ""
	}
	return jcloudDir
}

func CreateJlogDir(jcloudDir string) string {
	jlogDir := filepath.Join(jcloudDir, ".jlog")
	err := os.MkdirAll(jlogDir, 0750)
	if err != nil {
		return ""
	}
	return jlogDir
}

func CreateAnchorDir(jcloudDir string) string {
	anchorDir := filepath.Join(jcloudDir, ".anchor")
	err := os.MkdirAll(anchorDir, 0750)
	if err != nil {
		return ""
	}
	return anchorDir
}

func CreateJcloudFile(jcloudDir string) *os.File {
	file, err := os.OpenFile(filepath.Join(jcloudDir, ".jcloud"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil
	}
	return file
}

func CreateLogFile(jlogDir string) *os.File {
	file, err := os.OpenFile(filepath.Join(jlogDir, "jlog.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil
	}
	return file
}

func CreateJcookieFile(jcloudDir string) *os.File {
	file, err := os.OpenFile(filepath.Join(jcloudDir, ".jcookie"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil
	}
	return file
}

func CreateAnchorFile(anchorDir string) *os.File {
	file, err := os.OpenFile(filepath.Join(anchorDir, "anchor.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil
	}
	return file
}

func SetPaths() *Paths {
	homeDir := GetHome()
	jcloudDir := CreateJcloudDir(homeDir)
	jlogDir := CreateJlogDir(jcloudDir)
	anchorDir := CreateAnchorDir(jcloudDir)
	anchorLog := CreateAnchorFile(anchorDir)

	jcloudFile := CreateJcloudFile(jcloudDir)
	logFile := CreateLogFile(jlogDir)
	jcookieFile := CreateJcookieFile(jcloudDir)

	return &Paths{
		Home:       homeDir,
		JcloudFile: jcloudFile,
		JlogFile:   logFile,
		Jcookie:    jcookieFile,
		AnchorFile: anchorLog,
	}
}

func (p *Paths) Close() {
	err := p.JcloudFile.Close()
	if err != nil {
		return
	}
	err = p.Jcookie.Close()
	if err != nil {
		return
	}
	err = p.AnchorFile.Close()
	if err != nil {
		return
	}
	err = p.JlogFile.Close()
	if err != nil {
		return
	}
}
