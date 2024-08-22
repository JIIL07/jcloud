package cmd

import (
	"errors"
	"github.com/JIIL07/jcloud/internal/client/lib/logger"
	"log"

	"github.com/spf13/cobra"
)

var dropFlag bool

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a resource to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		if dropFlag {
			err := fctx.AddFileFromExplorer()
			if err != nil {
				log.Println(err)
			}
		} else {
			fctx.Logger.Error("Please use -d flag to drop a file from an opened explorer", slg.Err(errors.New("method not implemented")))
		}

	},
}

func init() {
	addCmd.Flags().BoolVarP(&dropFlag, "drop", "d", false, "Drop a file from an opened explorer")
	RootCmd.AddCommand(addCmd)
}
