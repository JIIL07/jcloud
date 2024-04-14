package main

import (
	"bufio"
	"os"
)

type Info struct {
	Fullname  string
	Id        int
	Filename  string
	Extension string
	Filesize  int
	Status    string
	Data      []byte
}

type Upgateinfo struct {
	fullNotation string
	name         string
	ext          string
}

type SearchInfo struct {
	fullNotation string
	name         string
	ext          string
}

type CreateFile struct {
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
var search SearchInfo
var update Upgateinfo
var createFile CreateFile

var reader = bufio.NewReader(os.Stdin)
var err error
