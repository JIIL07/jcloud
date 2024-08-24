package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "exit JcloudFile CLI",
	Long:  "exit JcloudFile CLI, remove local storage and write logs",
	Run: func(cmd *cobra.Command, args []string) {
		fs.StorageService.S.Close()
		filePath := "./local.db"
		err := os.Remove(filePath)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(exitCmd)
}
