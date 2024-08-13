package cmd

import (
	"github.com/JIIL07/jcloud/internal/client/requests"
	"log"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a resource to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		f := &requests.File{
			Filename:  args[0],
			Extension: "extension",
			Filesize:  1024,
			Status:    "status",
			Data:      []byte("data"),
		}
		err := requests.UploadFile(f)
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
