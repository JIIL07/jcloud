package cmd

import (
	"github.com/JIIL07/jcloud/internal/client/config"
	"github.com/JIIL07/jcloud/internal/client/storage"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the cloud",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.MustLoad()
		s, err := storage.InitDatabase(c)
		if err != nil {
			return err
		}
		err = s.CreateTable()
		if err != nil {
			return err
		}
		ctx.Storage = &s
		return nil

	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
