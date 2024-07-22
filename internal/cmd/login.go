<<<<<<< e11e69cc92d7c5a0aef5e57c3c44bea1b6154e12:cmd/login.go
package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	cloud "github.com/JIIL07/cloudFiles-manager/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		user, password := cloud.Login()

		encryptedUsername, err := cloud.Encrypt([]byte(user))

		if err != nil {
			log.Printf("error encrypting usernames: %v", err)
		}

		encryptedPassword, err := cloud.Encrypt([]byte(password))
		if err != nil {
			log.Printf("error encrypting password: %v", err)
		}

		viper.Set("username", base64.StdEncoding.EncodeToString(encryptedUsername))
		viper.Set("password", base64.StdEncoding.EncodeToString(encryptedPassword))

		saveConfig()

		err = pushlogin(base64.StdEncoding.EncodeToString(encryptedUsername), base64.StdEncoding.EncodeToString(encryptedPassword))
		if err != nil {
			log.Printf("error pushing login: %v", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}

func pushlogin(us, ps string) error {
	data := map[string]string{
		"username": us,
		"password": ps,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json encoding error: %v", err)
	}

	req, err := http.NewRequest("POST", "https://cloudfiles.up.railway.app/user", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request error: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
=======
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	cloud "github.com/JIIL07/cloudFiles-manager/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		username, password := cloud.Login()
		password = cloud.HashPassword(password)

		viper.Set("username", username)
		viper.Set("password", password)

		saveConfig()

		err := pushlogin(username, password)
		if err != nil {
			log.Printf("error pushing login: %v", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}

func pushlogin(us, ps string) error {
	data := map[string]string{
		"username": us,
		"password": ps,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json encoding error: %v", err)
	}

	req, err := http.NewRequest("POST", "/adduser", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request error: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
>>>>>>> Big file structure update, dockerfile does not work currently, some new features in code: new server system (not complete) and new logger:internal/cmd/login.go
