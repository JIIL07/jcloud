package main

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
	Text      string
}

type SearchInfo struct {
	fullNotation string
	name         string
	ext          string
}

var Statuses = []string{
	"Created",
	"Modified",
	"Deleted",
	"Renamed",
}
var info Info
var search SearchInfo
var reader = bufio.NewReader(os.Stdin)
