package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sync"
)

func allFiles(args []string, workerFunc func(string)) {
	targetDir, err := getTargetDir(args)
	if err != nil {
		cobra.CheckErr(err)
	}

	files := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go worker(files, &wg, workerFunc)
	}

	err = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if ignoreErrors {
				logVerbose("Error walking directory", "path", path, "error", err)
				return nil
			}
			return err
		}

		if excludeFile(path) {
			return nil
		}

		if !info.IsDir() {
			wg.Add(1)
			files <- path
		}

		return nil
	})

	close(files)
	wg.Wait()

	if err != nil && !ignoreErrors {
		cobra.CheckErr(err)
	}
}

func withInteractive(args []string, workerFunc func(string)) {
	for _, arg := range args {
		if confirmAction(arg) {
			workerFunc(arg)
		} else {
			logVerbose("Skipping", "file", arg)
		}
	}
}

func withWorkerPool(args []string, workerFunc func(string)) {
	files := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go worker(files, &wg, workerFunc)
	}

	for _, arg := range args {
		wg.Add(1)
		files <- arg
	}

	close(files)
	wg.Wait()
}

func worker(files chan string, wg *sync.WaitGroup, workerFunc func(string)) {
	for file := range files {
		workerFunc(file)
		wg.Done()
	}
}
