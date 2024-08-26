package cmd

import (
	"errors"
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/lib/cookies"
	jhash "github.com/JIIL07/jcloud/internal/client/lib/hash"
	"github.com/JIIL07/jcloud/internal/client/requests"
	"github.com/spf13/cobra"
	"os"
)

var loginCmd = &cobra.Command{
	Use:     "login [args]",
	Short:   "login to jcloud",
	Long:    "login to jcloud locally and store user credentials in .jcloud file, send email credentials to jcloud server to store in database",
	Example: "jcloud login [username] [email] [password]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("not enough arguments")
		}

		err := os.WriteFile(appCtx.PathsService.P.JcloudFile.Name(), []byte(args[0]+" "+args[1]+" "+jhash.Hash(args[2])), os.ModePerm)
		if err != nil {
			return err
		}

		u := &requests.UserData{
			Username: args[0],
			Email:    args[1],
			Password: jhash.Hash(args[2]),
		}
		resp, err := requests.Login(u)
		if err != nil {
			return err
		}

		rawCookies, err := cookies.Serialize(resp.Cookies())
		if err != nil {
			return err
		}

		err = cookies.WriteToFile(appCtx.PathsService.P.Jcookie.Name(), rawCookies)
		if err != nil {
			return err
		}

		appCtx.LoggerService.L.Info(fmt.Sprintf("new user %v logged in with email %v", args[0], args[1]))

		return nil
	},
}

func init() {
	//loginCmd.SetHelpFunc(customLoginHelpFunc)

	RootCmd.AddCommand(loginCmd)
}

func customLoginHelpFunc(cmd *cobra.Command, args []string) {
	helpMessage := fmt.Sprintf(`%s

Usage:
  jcloud add [args]

Example:
  jcloud login [username] [email] [password] (separated by spaces)

Use "jcloud [command] --help" for more information about a command.
`, loginCmd.Long)
	fmt.Print(helpMessage)
}
