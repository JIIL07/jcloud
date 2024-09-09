package cmd

import (
	"fmt"
	h "github.com/JIIL07/jcloud/internal/client/hints"
	"github.com/JIIL07/jcloud/internal/client/jc"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	dropFlag     bool
	allFlag      bool
	ignoreErrors bool
	dryRun       bool
	updateFlag   bool
	forceFlag    bool
	verboseFlag  bool
	interactive  bool
	excludeFiles []string
	targetDir    string
	hintsEnabled = true

	mutex      sync.Mutex
	numWorkers = runtime.NumCPU() + 2
)

var addCmd = &cobra.Command{
	Use:   "add [flags] [path] [file]...",
	Short: "Add files or directories to local storage",
	Long:  `Add one or more files or directories to local storage (SQLite) before uploading to the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

		if len(args) == 0 && !allFlag && !dropFlag && !dryRun && !updateFlag {
			if hintsEnabled {
				hintMessage := h.DisplayHint("add", h.EmptyPath, nil)
				if hintMessage != "" {
					cobra.WriteStringAndCheck(os.Stdout, hintMessage)
				}
			}
			return
		}

		if len(args) == 0 && !allFlag {
			if hintsEnabled {
				hintMessage := h.DisplayHint("add", h.AllFlagMissing, nil)
				if hintMessage != "" {
					cobra.WriteStringAndCheck(os.Stdout, hintMessage)
				}
			}
			return
		}

		if len(args) == 0 && (updateFlag || dryRun || allFlag) {
			args = append(args, ".")
		}

		//TODO:
		// Check for modified files when --update is used, if none exist, show hint
		//if updateFlag && !hasModifiedFiles(args) {
		//	if hintsEnabled {
		//		hintMessage := h.DisplayHint("add", h.NoModifiedFiles, nil)
		//		if hintMessage != "" {
		//			cobra.WriteStringAndCheck(os.Stdout, hintMessage)
		//		}
		//	}
		//	return
		//}

		if allFlag || (len(args) == 1 && args[0] == ".") {
			addAll(args)
			return
		}

		if interactive {
			withInteractive(args)
		} else {
			withWorkerPool(args)
		}

		logVerbose("Command execution time", "duration", time.Since(startTime))
		appCtx.LoggerService.L.Info("Command execution time", "duration", time.Since(startTime))
	},
}

func workerPool(files chan string, wg *sync.WaitGroup) {
	for file := range files {
		add(file)
		wg.Done()
	}
}

func withInteractive(args []string) {
	for _, arg := range args {
		add(arg)
	}
}

func withWorkerPool(args []string) {
	files := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go workerPool(files, &wg)
	}

	for _, arg := range args {
		wg.Add(1)
		files <- arg
	}

	close(files)
	wg.Wait()
}

func add(arg string) {
	info, err := os.Stat(arg)
	if os.IsNotExist(err) {
		hintMessage := h.DisplayHint("add", h.NoFilesMatched, nil)
		if hintMessage != "" {
			cobra.WriteStringAndCheck(os.Stdout, hintMessage)
		}
		logVerbose("File or directory does not exist", "path", arg)
		if !ignoreErrors {
			return
		}
	}

	if err != nil {
		logVerbose("Error accessing file or directory", "path", arg, "error", err)
		cobra.CheckErr(err)
	}

	if info.IsDir() {
		addAll([]string{arg})
	} else {
		if interactive {
			var response string
			fmt.Printf("Add file %s? (y/n): ", arg)
			_, err := fmt.Scanln(&response)
			if err != nil || (response != "y" && response != "Y") {
				logVerbose("Skipping file", "file", arg)
				return
			}
		}

		if dryRun {
			cobra.WriteStringAndCheck(os.Stdout, fmt.Sprintf("Would add file %s\n", arg))
		} else {
			mutex.Lock()
			defer mutex.Unlock()

			logVerbose("Adding file", "file", arg)
			err = jc.AddFileFromPath(appCtx.FileService, arg)
			if err != nil && !ignoreErrors {
				cobra.CheckErr(fmt.Errorf("failed to add file: %w", err))
			}
		}
	}
}

func addAll(args []string) {
	var err error
	targetDir, err = getTargetDir(args)
	if err != nil {
		cobra.CheckErr(err)
	}

	files := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go workerPool(files, &wg)
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

		if info.IsDir() {
			return nil
		}

		logVerbose("Processing file", "file", path)
		wg.Add(1)
		files <- path
		return nil
	})

	close(files)
	wg.Wait()

	if err != nil && !ignoreErrors {
		cobra.CheckErr(err)
	}
}

func excludeFile(path string) bool {
	for _, exclude := range excludeFiles {
		absExclude, err := filepath.Abs(exclude)
		if err != nil {
			logVerbose("Error resolving exclude path", "exclude", exclude, "error", err)
			continue
		}

		if strings.HasPrefix(path, absExclude) {
			logVerbose("Excluding path", "file", path)
			return true
		}
	}
	return false
}

func getTargetDir(args []string) (string, error) {
	if args != nil {
		return filepath.Abs(args[0])
	}
	return os.Getwd()
}

func logVerbose(message string, keysAndValues ...interface{}) {
	if verboseFlag {
		cobra.WriteStringAndCheck(os.Stdout, fmt.Sprintf("[%s] %s - %v\n", time.Now().Format(time.RFC3339), message, keysAndValues))
	}
}

func init() {
	addCmd.Flags().BoolVarP(&dropFlag, "drop", "d", false, "Drop a file from an opened explorer")
	addCmd.Flags().BoolVarP(&allFlag, "all", "A", false, "Add all files in the current directory (recursive)")
	addCmd.Flags().BoolVarP(&ignoreErrors, "ignore-errors", "I", false, "Ignore errors and continue adding files")
	addCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Show files that would be added, without adding them")
	addCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "Only add modified files, skip untracked ones")
	addCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force adding ignored files")
	addCmd.Flags().BoolVarP(&verboseFlag, "verbose", "V", false, "Show detailed logs during file addition")
	addCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactively choose files to add")
	addCmd.Flags().StringSliceVarP(&excludeFiles, "exclude", "e", []string{}, "Exclude specific files or directories")

	addCmd.Flags().BoolVar(&hintsEnabled, "advice", true, "Enable or disable hints when nothing is specified")

	addCmd.Args = cobra.MinimumNArgs(0)

	RootCmd.AddCommand(addCmd)
}
