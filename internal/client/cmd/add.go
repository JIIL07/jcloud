package cmd

import (
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
			log.Println("no -d --drop usage")
		}

	},
}

func init() {
	addCmd.Flags().BoolVarP(&dropFlag, "drop", "d", false, "Drop a file from an opened explorer")
	RootCmd.AddCommand(addCmd)
}
