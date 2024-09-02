package cmd

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/jc"
	"github.com/JIIL07/jcloud/pkg/log"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print list of files",
	Long:  "Print list of files from local storage",
	Run: func(cmd *cobra.Command, args []string) {
		items, err := jc.ListFiles(appCtx.FileService)
		if err != nil {
			appCtx.LoggerService.L.Error("error listing files", jlog.Err(err))
			cobra.CheckErr(err)
		}
		for _, item := range items {
			cobra.WriteStringAndCheck(os.Stdout, fmt.Sprintf("- %v\n", item))
			//TODO: implement good file output
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

}
