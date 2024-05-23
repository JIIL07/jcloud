package file

import "regexp"

var Statuses = []string{"Created", "Has data in", "Renamed", "Deleted"}

var isValidTableName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.]*$`).MatchString
var isValidDBName = regexp.MustCompile(`^[^<>:"/\\|?*]+$`).MatchString
