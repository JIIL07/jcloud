package cmd

import (
	"github.com/JIIL07/cloudFiles-manager/internal/client/requests"
	"github.com/spf13/cobra"
	"log"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		u := &requests.UserData{
			Username: args[0],
			Password: "password",
			Email:    "email3@gmail.com",
		}
		err := requests.Login(u)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
