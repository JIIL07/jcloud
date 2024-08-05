package cmd

import (
	"fmt"
	"log"
	"os"

	cloud "github.com/JIIL07/cloudFiles-manager/internal/client"
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

		_, _ = username, password
	}
}
