package cmd

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/jc"
	"github.com/JIIL07/jcloud/pkg/log"
	"github.com/spf13/cobra"
	"os"
)

var dataFlag bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print list of files",
	Long:  "Print list of files from local storage",
	Run: func(cmd *cobra.Command, args []string) {
		items, err := jc.ListFiles(a.File)
		if err != nil {
			a.Logger.L.Error("error listing files", jlog.Err(err))
			cobra.CheckErr(err)
		}
		for _, item := range items {
			cobra.WriteStringAndCheck(os.Stdout, fmt.Sprintf("ID: %d\n", item.ID))
			cobra.WriteStringAndCheck(os.Stdout, fmt.Sprintf("File: %s.%s\n", item.Meta.Name, item.Meta.Extension))
			cobra.WriteStringAndCheck(os.Stdout, fmt.Sprintf("Size: %d bytes\n", item.Meta.Size))
			if dataFlag {
				cobra.WriteStringAndCheck(os.Stdout, fmt.Sprintf("File content: %s\n\n", string(item.Data)))
			} else {
				cobra.WriteStringAndCheck(os.Stdout, fmt.Sprintln())
			}
		}
	},
}

func init() {
	listCmd.Flags().BoolVarP(&dataFlag, "data", "d", false, "Show file content")
	RootCmd.AddCommand(listCmd)
}
