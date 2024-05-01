module main

go 1.22.2

require project/file v0.0.1

require (
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
)

replace project/file => ../file
