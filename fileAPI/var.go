package file

import (
	"bufio"
	"os"
)

type Info struct {
	Id        int
	Filename  string
	Extension string
	Filesize  int
	Status    string
	Data      []byte

	Fullname string
}

type TempInfo struct {
	fullNotation string
	name         string
	ext          string
}

var Statuses = []string{
	"Created",
	"Has data in",
	"Renamed",
}
var info Info
var temp TempInfo

var reader = bufio.NewReader(os.Stdin)
var err error

var SomeVariable string
