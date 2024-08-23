package cmd

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/lib/logger"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print list of files",
	Long:  "Print list of files from local storage",
	Run: func(cmd *cobra.Command, args []string) {
		items, err := fctx.ListFiles()
		if err != nil {
			logger.Error("error listing files", slg.Err(err))
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
