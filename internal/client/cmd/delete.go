package cmd

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/jc"
	"github.com/spf13/cobra"
	"os"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [flags] [filename]...",
	Short: "Delete file or directory records from local database",
	Long:  "Delete one or more file or directory records from the local SQLite database without affecting the server storage.",
	Run: func(cmd *cobra.Command, args []string) {
		if allFlag {
			confirmAction("all files")
			err := jc.DeleteAllFiles(a.File)
			if err != nil {
				cobra.CheckErr(fmt.Errorf("failed to delete all files from database: %w", err))
			}
			return
		}

		if len(args) == 0 {
			cobra.CheckErr(fmt.Errorf("no files specified"))
		}

		if interactive {
			withInteractive(args, deleteFile)
		} else {
			withWorkerPool(args, deleteFile)
		}
	},
}

func deleteFile(arg string) {
	if arg == "" {
		cobra.CheckErr(fmt.Errorf("no file specified"))
		return
	}

	if dryRun {
		cobra.WriteStringAndCheck(os.Stdout, fmt.Sprintf("Would delete record for file %s from database\n", arg))
	} else {
		logVerbose("Deleting file record", "file", arg)
		a.File.F.Meta.Name = arg
		err := jc.DeleteFile(a.File)
		if err != nil && !ignoreErrors {
			cobra.CheckErr(fmt.Errorf("failed to delete file record: %w", err))
		}
	}
}

func init() {
	deleteCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Delete all files")
	deleteCmd.Flags().BoolVarP(&ignoreErrors, "ignore-errors", "I", false, "Ignore errors and continue deleting files")
	deleteCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Show files that would be deleted, without deleting them")
	deleteCmd.Flags().BoolVarP(&verboseFlag, "verbose", "V", false, "Show detailed logs during deletion")
	deleteCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactively choose files to delete")

	RootCmd.AddCommand(deleteCmd)
}
