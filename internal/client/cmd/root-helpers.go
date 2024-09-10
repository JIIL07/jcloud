package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func confirmAction(arg string) bool {
	defaultResponse := "y"
	fmt.Printf("Are you sure you want to process %s? (Y/n): ", arg)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.TrimSpace(scanner.Text())

	if response == "" {
		response = defaultResponse
	}

	return response == "y" || response == "Y"
}

func excludeFile(path string) bool {
	for _, exclude := range excludeFiles {
		absExclude, err := filepath.Abs(exclude)
		if err != nil {
			logVerbose("Error resolving exclude path", "exclude", exclude, "error", err)
			continue
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			logVerbose("Error resolving file path", "file", path, "error", err)
			continue
		}

		matched, err := filepath.Match(exclude, filepath.Base(path))
		if err == nil && matched {
			logVerbose("Excluding path by pattern", "file", path)
			return true
		}

		rel, err := filepath.Rel(absExclude, absPath)
		if err != nil {
			logVerbose("Error getting relative path", "base", absExclude, "target", absPath, "error", err)
			continue
		}

		if !strings.HasPrefix(rel, "..") {
			logVerbose("Excluding path", "file", absPath)
			return true
		}
	}
	return false
}

func getTargetDir(args []string) (string, error) {
	if len(args) > 0 {
		return filepath.Abs(args[0])
	}
	return os.Getwd()
}

func logVerbose(message string, keysAndValues ...interface{}) {
	if verboseFlag {
		fmt.Printf("[%s] %s - %v\n", time.Now().Format(time.RFC3339), message, keysAndValues)
	}
}
