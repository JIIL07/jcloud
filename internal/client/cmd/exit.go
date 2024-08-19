package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// exitCmd represents the exit command
var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Exit the cloud CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fctx.Storage.Close()
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
