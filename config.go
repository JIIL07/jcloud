package cloudfiles

import "regexp"

var Statuses = []string{"Created", "Has data in", "Renamed", "Deleted"}

var isValidTableName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.]*$`).MatchString
