package file

import "regexp"

// Statuses defines possible statuses of a file.
var Statuses = []string{"Created", "Has data in", "Renamed", "Deleted"}

// Regular expressions for validation
const (
	tableNamePattern = `^[a-zA-Z_][a-zA-Z0-9_.]*$`
	dbNamePattern    = `^[^<>:"/\\|?*]+$`
)

// isValidTableName checks if the provided string is a valid table name.
var isValidTableName = regexp.MustCompile(tableNamePattern).MatchString

// isValidDBName checks if the provided string is a valid database name.
var isValidDBName = regexp.MustCompile(dbNamePattern).MatchString
