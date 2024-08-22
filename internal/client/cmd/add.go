package cmd

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/lib/logger"
	"os"

	"github.com/spf13/cobra"
)

var (
	dropFlag bool
	filepath string
	allFiles bool
)

var addCmd = &cobra.Command{
	Use:   "add [flags] | [filename]",
	Short: "Add a file",
	Long:  "Add a file to local storage (SQLite) before uploading it to the server",
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case dropFlag:
			err := fctx.AddFileFromExplorer()
			if err != nil {
				fctx.Logger.Error("error adding file via drop-down", slg.Err(err))
				cobra.CheckErr(err)
			}
		case filepath != "":
			// TODO: Add file from path
		case allFiles:
			// TODO: Add all files from current directory
		default:
			cobra.WriteStringAndCheck(os.Stdout, "specify a flag or use 'jcloud add --help' for more information\n")
		}

	},
}

func init() {
	addCmd.Flags().BoolVarP(&dropFlag, "drop", "d", false, "Drop a file from an opened explorer")
	addCmd.Flags().StringVarP(&filepath, "path", "p", "", "Add a file from path")
	addCmd.Flags().BoolVarP(&allFiles, "all", "a", false, "Add all files from current directory")

	addCmd.SetHelpFunc(customAddHelpFunc)

	RootCmd.AddCommand(addCmd)
}

func customAddHelpFunc(cmd *cobra.Command, args []string) {
	helpMessage := fmt.Sprintf(`%s

Usage:
  jcloud add [flags]

Flags:
  -d, --drop             Drop a file from an opened explorer
  -p, --path [filepath]  Add a file from path
  -a, --all              Add all files from current directory

Use "jcloud [command] --help" for more information about a command.
`, addCmd.Long)
	fmt.Print(helpMessage)
}
