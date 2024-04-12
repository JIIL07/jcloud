package main

type Info struct {
	Id        int
	Filename  string
	Extension string
	Filesize  int
	Status    string
	Text      string
}

var Statuses = []string{
	"Created",
	"Modified",
	"Deleted",
	"Renamed",
}
var info Info
