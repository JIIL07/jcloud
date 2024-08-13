package cmd

import (
	"github.com/JIIL07/jcloud/internal/client/requests"
	"github.com/spf13/cobra"
	"log"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := requests.GetFiles()
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
