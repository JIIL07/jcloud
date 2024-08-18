package cmd

import (
	"github.com/JIIL07/jcloud/internal/client/models"
	cloud "github.com/JIIL07/jcloud/internal/client/operations"
	"github.com/spf13/cobra"
)

var ctx *cloud.FileContext

var RootCmd = &cobra.Command{
	Use:   `cloud`,
	Short: `Cloud is a cloud file manager CLI`,
	Long: `Cloud is a cloud file manager CLI that provides various commands to manage files in the cloud.
It supports commands like init to initialize the cloud, add to add files, and exit to exit the CLI.`,
}

func init() {
	ctx = &cloud.FileContext{
		Info: &models.File{},
	}
	RootCmd.PersistentFlags().BoolP("help", "h", false, "Help")
}
