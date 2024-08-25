package cmd

import (
	"github.com/JIIL07/jcloud/internal/client/jc"
	slg "github.com/JIIL07/jcloud/internal/client/lib/logger"
	"github.com/spf13/cobra"
)

var allFilesD bool

var deleteCmd = &cobra.Command{
	Use:   "delete [flags] | [filename]",
	Short: "Delete file",
	Long:  "Delete file from local storage do not collide with server storage",
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case len(args) > 0:
			appCtx.FileService.F.Metadata.Name = args[0]
			err := jc.DeleteFile(appCtx.FileService)
			if err != nil {
				appCtx.LoggerService.L.Error("error deleting file", slg.Err(err))
				cobra.CheckErr(err)
			}
		case allFilesD:
			err := jc.DeleteAllFiles(appCtx.FileService)
			if err != nil {
				appCtx.LoggerService.L.Error("error deleting all files", slg.Err(err))
				cobra.CheckErr(err)
			}
		}
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&allFilesD, "all", "a", false, "Delete all files from local storage")
	RootCmd.AddCommand(deleteCmd)
}
