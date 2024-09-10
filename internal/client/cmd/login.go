package cmd

import (
	"errors"
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/requests"
	"github.com/JIIL07/jcloud/pkg/cookies"
	jhash "github.com/JIIL07/jcloud/pkg/hash"
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

		err := os.WriteFile(a.Paths.P.JcloudFile.Name(), []byte(args[0]+" "+args[1]+" "+jhash.Hash(args[2])), 0600)
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

		err = cookies.WriteToFile(a.Paths.P.Jcookie.Name(), rawCookies)
		if err != nil {
			return err
		}

		a.Logger.L.Info(fmt.Sprintf("new user %v logged in with email %v", args[0], args[1]))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
