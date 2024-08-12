package cmd

import (
	"github.com/JIIL07/cloudFiles-manager/internal/client/config"
	"github.com/JIIL07/cloudFiles-manager/internal/client/storage"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the cloud",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.MustLoad()
		sqlite, err := storage.InitDatabase(c)
		if err != nil {
			return err
		}
		err = sqlite.CreateTable("files")
		if err != nil {
			return err
		}
		return nil

	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
