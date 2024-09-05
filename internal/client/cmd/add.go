package cmd

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/jc"
	"github.com/JIIL07/jcloud/pkg/log"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	dropFlag     bool
	targetDir    string
	allFlag      bool
	ignoreErrors bool
	dryRun       bool
	updateFlag   bool
)

var addCmd = &cobra.Command{
	Use:   "add [flags] [file or directory]...",
	Short: "Add files or directories to local storage",
	Long:  `Add one or more files or directories to local storage (SQLite) before uploading to the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()
		if len(args) == 0 && (updateFlag || dryRun || allFlag) {
			args = append(args, ".")
		}

		if allFlag || (len(args) == 1 && args[0] == ".") {
			var err error
			targetDir, err = os.Getwd()
			if err != nil {
				appCtx.LoggerService.L.Error("error getting current directory", jlog.Err(err))
				cobra.CheckErr(err)
			}

			err = jc.AddFilesFromDir(appCtx.FileService, targetDir)
			if err != nil {
				if !ignoreErrors {
					cobra.CheckErr(fmt.Errorf("failed to add files from directory: %w", err))
				} else {
					appCtx.LoggerService.L.Warn("error adding files", jlog.Err(err))
				}
			}
			return
		}

		if dropFlag {
			err := jc.AddFileFromExplorer(appCtx.FileService)
			if err != nil {
				appCtx.LoggerService.L.Error("error adding file via drop-down", jlog.Err(err))
				cobra.CheckErr(err)
			}
			return
		}

		if len(args) == 0 && !dropFlag {
			cobra.CheckErr(fmt.Errorf("no files or directories specified"))
		}

		for _, arg := range args {
			info, err := os.Stat(arg)
			if os.IsNotExist(err) {
				appCtx.LoggerService.L.Warn("file or directory does not exist", jlog.Err(err), "path", arg)
				if !ignoreErrors {
					continue
				}
			}
			if err != nil {
				appCtx.LoggerService.L.Error("error accessing file or directory", jlog.Err(err), "path", arg)
				cobra.CheckErr(err)
			}

			if updateFlag && !info.ModTime().After(time.Now()) {
				appCtx.LoggerService.L.Info("skipping unmodified file", "file", arg)
				continue
			}

			if info.IsDir() {
				err = jc.AddFilesFromDir(appCtx.FileService, arg)
				if err != nil {
					if !ignoreErrors {
						cobra.CheckErr(fmt.Errorf("failed to add files from directory: %w", err))
					} else {
						appCtx.LoggerService.L.Warn("error adding files", jlog.Err(err))
					}
				}
			} else {
				if dryRun {
					appCtx.LoggerService.L.Info("would add file", "file", arg)
				} else {
					err = jc.AddFileFromPath(appCtx.FileService, arg)
					if err != nil {
						if !ignoreErrors {
							cobra.CheckErr(fmt.Errorf("failed to add file: %w", err))
						} else {
							appCtx.LoggerService.L.Warn("error adding file", jlog.Err(err), "file", arg)
						}
					}
				}
			}
		}
		elapsedTime := time.Since(startTime)
		appCtx.LoggerService.L.Info("Command execution time", "duration", elapsedTime)
	},
}

func init() {
	addCmd.Flags().BoolVarP(&dropFlag, "drop", "d", false, "Drop a file from an opened explorer")
	addCmd.Flags().BoolVarP(&allFlag, "all", "A", false, "Add all files in the current directory (recursive)")
	addCmd.Flags().BoolVarP(&ignoreErrors, "ignore-errors", "i", false, "Ignore errors and continue adding files")
	addCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Show files that would be added, without adding them")
	addCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "Only add modified files, skip untracked ones")
	addCmd.Flags().StringVarP(&targetDir, "path", "p", "", "Specify the directory to add all files from")

	addCmd.Args = cobra.MinimumNArgs(0)

	RootCmd.AddCommand(addCmd)
}
