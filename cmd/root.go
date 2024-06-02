package cmd

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	cloud "github.com/JIIL07/cloudFiles-manager/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ctx *cloud.FileContext

var RootCmd = &cobra.Command{
	Use:   `cloud`,
	Short: `Cloud is a cloud file manager CLI`,
	Long: `Cloud is a cloud file manager CLI that provides various commands to manage files in the cloud.
It supports commands like init to initialize the cloud, add to add files, and exit to exit the CLI.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "login" {
			if _, err := os.Stat(fmt.Sprintf("%s/.cloud_init_config.json", os.Getenv("TEMP"))); err != nil {
				log.Println("Use cloud login to create cloud configuration")
			}
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolP("help", "h", false, "Help")

	viper.SetConfigName(".cloud_init_config")
	viper.SetConfigType("json")
	viper.AddConfigPath(os.Getenv("TEMP"))

	cloud.Generator(32)
}

func saveConfig() {
	if err := viper.WriteConfigAs(fmt.Sprintf("%s/.cloud_init_config.json", os.Getenv("TEMP"))); err != nil {
		log.Println("Error writing config file:", err)
	}
}

func loadUserPass() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Println("Error reading config file:", err)
		}
	} else {
		username := viper.GetString("username")
		password := viper.GetString("password")

		decodedUsername, err := base64.StdEncoding.DecodeString(username)
		if err != nil {
			log.Println("Error decoding username:", err)
		}
		decodedPassword, err := base64.StdEncoding.DecodeString(password)
		if err != nil {
			log.Println("Error decoding password:", err)
		}

		userByte, err := cloud.Decrypt(decodedUsername)
		if err != nil {
			log.Println("Error decrypting username:", err)
		}
		passByte, err := cloud.Decrypt(decodedPassword)
		if err != nil {
			log.Println("Error decrypting password:", err)
		}
		fmt.Println("Username:", string(userByte))
		fmt.Println("Password:", string(passByte))
	}
}
