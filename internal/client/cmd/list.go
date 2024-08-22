package cmd

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/lib/logger"
	"os"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print list of files",
	Long:  "Print list of files from local storage",
	Run: func(cmd *cobra.Command, args []string) {
		items, err := fctx.ListFiles()
		if err != nil {
			fctx.Logger.Error("error listing files", slg.Err(err))
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
