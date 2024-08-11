package cmd

import (
	"github.com/JIIL07/cloudFiles-manager/internal/client/requests"
	"log"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a resource to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		f := &requests.File{}
		err := requests.UploadFile(f)
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
