package cmd

import (
	"fmt"
	h "github.com/JIIL07/jcloud/internal/client/hints"
	"github.com/JIIL07/jcloud/internal/client/jc"
	"github.com/JIIL07/jcloud/pkg/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	dropFlag     bool
	targetDir    string
	allFlag      bool
	ignoreErrors bool
	dryRun       bool
	updateFlag   bool
	forceFlag    bool
	verboseFlag  bool
	interactive  bool
	excludeFiles []string
	hintsEnabled bool = true
)

var addCmd = &cobra.Command{
	Use:   "add [flags] [file or directory]...",
	Short: "Add files or directories to local storage",
	Long:  `Add one or more files or directories to local storage (SQLite) before uploading to the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

		if len(args) == 0 && !allFlag && !dropFlag && !dryRun && !updateFlag {
			if hintsEnabled {
				hintMessage := h.DisplayHint("add", h.EmptyPath, c)
				if hintMessage != "" {
					cobra.WriteStringAndCheck(os.Stdout, hintMessage)
				}
			}
			return
		}

		if len(args) == 0 && (updateFlag || dryRun || allFlag) {
			args = append(args, ".")
		}

		if allFlag || (len(args) == 1 && args[0] == ".") {
			handleAddAllFiles()
			return
		}

		// Handle the --drop flag
		if dropFlag {
			handleDropAdd()
			return
		}

		// Handle specific file/directory arguments
		var wg sync.WaitGroup
		for _, arg := range args {
			wg.Add(1)
			go func(arg string) {
				defer wg.Done()
				handleAddArg(arg)
			}(arg)
		}
		wg.Wait()

		// Log command execution time
		appCtx.LoggerService.L.Info("Command execution time", "duration", time.Since(startTime))
	},
}

func handleAddAllFiles() {
	var err error
	targetDir, err = os.Getwd()
	if err != nil {
		appCtx.LoggerService.L.Error("error getting current directory", jlog.Err(err))
		cobra.CheckErr(err)
	}

	err = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if ignoreErrors {
				appCtx.LoggerService.L.Warn("error walking directory", jlog.Err(err), "path", path)
				return nil
			}
			return err
		}

		// Exclude files based on pattern
		for _, exclude := range excludeFiles {
			if match, _ := filepath.Match(exclude, path); match {
				appCtx.LoggerService.L.Info("excluding file", "file", path)
				return nil
			}
		}

		// Process directories and files
		if info.IsDir() {
			if verboseFlag {
				appCtx.LoggerService.L.Info("adding directory", "dir", path)
			}
			err = jc.AddFilesFromDir(appCtx.FileService, path)
		} else {
			if dryRun {
				appCtx.LoggerService.L.Info("would add file", "file", path)
			} else {
				if verboseFlag {
					appCtx.LoggerService.L.Info("adding file", "file", path)
				}
				err = jc.AddFileFromPath(appCtx.FileService, path)
			}
		}
		if err != nil && !ignoreErrors {
			return fmt.Errorf("failed to add files from directory: %w", err)
		}
		return nil
	})

	if err != nil && !ignoreErrors {
		cobra.CheckErr(err)
	}
}

func handleDropAdd() {
	err := jc.AddFileFromExplorer(appCtx.FileService)
	if err != nil {
		appCtx.LoggerService.L.Error("error adding file via drop-down", jlog.Err(err))
		cobra.CheckErr(err)
	}
}

func handleAddArg(arg string) {
	info, err := os.Stat(arg)
	if os.IsNotExist(err) {
		// Display hint if the file or directory doesn't exist
		hintMessage := h.DisplayHint("add", h.NoFilesMatched, c)
		if hintMessage != "" {
			cobra.WriteStringAndCheck(os.Stdout, hintMessage)
		}
		appCtx.LoggerService.L.Warn("file or directory does not exist", jlog.Err(err), "path", arg)
		if !ignoreErrors {
			return
		}
	}

	if err != nil {
		appCtx.LoggerService.L.Error("error accessing file or directory", jlog.Err(err), "path", arg)
		cobra.CheckErr(err)
	}

	// Check for --update flag and only add modified files
	if updateFlag {
		// TODO: Replace with actual logic to get the last modified time from storage
		lastModTime := time.Now().Add(-1 * time.Hour)
		if !info.ModTime().After(lastModTime) {
			// Display hint for unmodified files
			hintMessage := h.DisplayHint("add", h.NoModifiedFiles, c)
			if hintMessage != "" {
				cobra.WriteStringAndCheck(os.Stdout, hintMessage)
			}
			appCtx.LoggerService.L.Info("skipping unmodified file", "file", arg)
			return
		}
	}

	// Handle directories
	if info.IsDir() {
		handleAddAllFiles()
	} else {
		if dryRun {
			appCtx.LoggerService.L.Info("would add file", "file", arg)
		} else {
			err = jc.AddFileFromPath(appCtx.FileService, arg)
			if err != nil && !ignoreErrors {
				cobra.CheckErr(fmt.Errorf("failed to add file: %w", err))
			}
		}
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
