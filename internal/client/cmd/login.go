package cmd

import (
	"errors"
	"fmt"
	jhash "github.com/JIIL07/jcloud/internal/client/lib/hash"
	"github.com/spf13/cobra"
	"os"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("not enough arguments")
		}

		err := os.WriteFile(fctx.Local.Jcloud, []byte(args[0]+" "+args[1]+" "+jhash.Hash(args[2])), os.ModePerm)
		if err != nil {
			return err
		}

		//u := &requests.UserData{
		//	Username: args[0],
		//	Email:    args[1],
		//	Password: jhash.Hash(args[2]),
		//}
		//err = requests.Login(u)
		//if err != nil {
		//	return err
		//}

		fctx.Logger.Info(fmt.Sprintf("new user %v logged in with email %v", args[0], args[1]))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
