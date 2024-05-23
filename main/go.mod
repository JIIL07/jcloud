module main

go 1.22.2

require project/file v0.0.1

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/sys v0.4.0 // indirect
)

replace project/file => ../file
