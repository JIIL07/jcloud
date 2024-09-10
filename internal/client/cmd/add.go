package cmd

import (
	"fmt"
	h "github.com/JIIL07/jcloud/internal/client/hints"
	"github.com/JIIL07/jcloud/internal/client/jc"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	dropFlag      bool
	updateFlag    bool
	forceFlag     bool
	interactiveV2 bool
	hintsEnabled  = true

	mutex      sync.Mutex
	numWorkers = runtime.NumCPU() + 2
)

var addCmd = &cobra.Command{
	Use:   "add [flags] [path] [file]...",
	Short: "Add files or directories to local storage",
	Long:  `Add one or more files or directories to local storage (SQLite) before uploading to the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()
		if interactiveV2 {
			withInteractive_v2(args)
		}

		if len(args) == 0 && !allFlag && !dropFlag && !dryRun && !updateFlag {
			if hintsEnabled {
				hintMessage := h.DisplayHint("addFile", h.EmptyPath, c)
				if hintMessage != "" {
					cobra.WriteStringAndCheck(os.Stdout, hintMessage)
				}
			}
			return
		}

		if len(args) == 0 && !allFlag {
			if hintsEnabled {
				hintMessage := h.DisplayHint("addFile", h.AllFlagMissing, c)
				if hintMessage != "" {
					cobra.WriteStringAndCheck(os.Stdout, hintMessage)
				}
			}
			return
		}

		if dropFlag {
			err := jc.AddFileFromExplorer(a.File)
			if err != nil {
				cobra.CheckErr(err)
			}
		}

		if allFlag || (len(args) == 1 && args[0] == ".") {
			allFiles(args, addFile)
		} else if interactive {
			withInteractive(args, addFile)
		} else {
			withWorkerPool(args, addFile)
		}

		logVerbose("Command execution time", "duration", time.Since(startTime))
		a.Logger.L.Info("Command execution time", "duration", time.Since(startTime))
	},
}

func addFile(arg string) {
	info, err := os.Stat(arg)
	if os.IsNotExist(err) {
		hintMessage := h.DisplayHint("addFile", h.NoFilesMatched, c)
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
		allFiles([]string{arg}, addFile)
	} else {
		if dryRun {
			cobra.WriteStringAndCheck(os.Stdout, fmt.Sprintf("Would addFile file %s\n", arg))
		} else {
			mutex.Lock()
			defer mutex.Unlock()

			logVerbose("Adding file", "file", arg)
			err = jc.AddFileFromPath(a.File, arg)
			if err != nil && !ignoreErrors {
				cobra.CheckErr(fmt.Errorf("failed to addFile file: %w", err))
			}
		}
	}
}

func init() {
	addCmd.Flags().BoolVarP(&dropFlag, "drop", "d", false, "Drop a file from an opened explorer")
	addCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Add all files in the current directory (recursive)")
	addCmd.Flags().BoolVarP(&ignoreErrors, "ignore-errors", "I", false, "Ignore errors and continue adding files")
	addCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Show files that would be added, without adding them")
	addCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "Only addFile modified files, skip untracked ones")
	addCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force adding ignored files")
	addCmd.Flags().BoolVarP(&verboseFlag, "verbose", "V", false, "Show detailed logs during file addition")
	addCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactively choose files to addFile")
	addCmd.Flags().BoolVar(&interactiveV2, "v2", false, "Interactively choose files to addFile")
	addCmd.Flags().StringSliceVarP(&excludeFiles, "exclude", "e", []string{}, "Exclude specific files or directories")

	addCmd.Flags().BoolVar(&hintsEnabled, "advice", true, "Enable or disable hints when nothing is specified")

	addCmd.Args = cobra.MinimumNArgs(0)

	RootCmd.AddCommand(addCmd)
}
