package cmd

import (
	cloud "github.com/JIIL07/cloudFiles-manager/client"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the cloud",
	RunE: func(cmd *cobra.Command, args []string) error {
		sqlite := &cloud.SQLiteDB{}
		db, err := sqlite.PrepareLocalDB()
		if err != nil {
			return err
		}
		ctx = cloud.NewFileContext(db)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
