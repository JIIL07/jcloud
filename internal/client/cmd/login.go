package cmd

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/requests"
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
		user := make([]string, 3)
		fmt.Print("Username: ")
		_, err := fmt.Scan(&user[0])
		if err != nil {
			return err
		}

		fmt.Print("Email: ")
		_, err = fmt.Scan(&user[1])
		if err != nil {
			return err
		}

		fmt.Print("Password: ")
		_, err = fmt.Scan(&user[2])
		if err != nil {
			return err
		}

		err = os.WriteFile(a.Paths.P.JcloudFile.Name(), []byte(user[0]+" "+user[1]+" "+jhash.Hash([]byte(user[2]))), 0600)
		if err != nil {
			return err
		}

		u := &requests.UserData{
			Username: user[0],
			Email:    user[1],
			Password: jhash.Hash([]byte(user[2])),
		}
		resp, err := requests.Login(u)
		if err != nil {
			return err
		}

		rawCookies, err := requests.Serialize(resp.Cookies())
		if err != nil {
			return err
		}

		err = requests.WriteToFile(a.Paths.P.Jcookie.Name(), rawCookies)
		if err != nil {
			return err
		}

		a.Logger.L.Info(fmt.Sprintf("new user %v logged in with email %v", user[0], user[1]))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
