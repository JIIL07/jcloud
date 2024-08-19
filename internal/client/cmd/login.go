package cmd

import (
	"errors"
	jhash "github.com/JIIL07/jcloud/internal/client/lib/hash"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var data = []byte(`syntax = "proto3";

message FileMetadata {
  string filename = 1;
  string extension = 2;
  int32 filesize = 3;
}

message File {
  int32 id = 1;
  FileMetadata metadata = 2;
  string status = 3;
  bytes data = 4;
}`)

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
		//u := &requests.UserData{
		//	Username: args[0],
		//	Password: "password",
		//	Email:    "email3@gmail.com",
		//}
		//err := requests.Login(u)
		//if err != nil {
		//	log.Fatal(err)
		//}
		if len(args) < 3 {
			return errors.New("not enough arguments")
		}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		jcloudDir := filepath.Join(homeDir, ".jcloud")
		err = os.MkdirAll(jcloudDir, os.ModePerm)
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(jcloudDir, ".jcloud"), []byte(args[0]+" "+args[1]+" "+jhash.Hash(args[2])), os.ModePerm)
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(jcloudDir, ".jcloud.proto"), []byte(data), os.ModePerm)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
