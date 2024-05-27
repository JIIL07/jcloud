package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	cloud "github.com/JIIL07/cloudFiles-manager/client"
	"github.com/spf13/cobra"
)

var ctx *cloud.FileContext

var RootCmd = &cobra.Command{
	Use:   `cloud`,
	Short: `Cloud is a cloud file manager CLI`,
	Long: `Cloud is a cloud file manager CLI that provides various commands to manage files in the cloud.
It supports commands like init to initialize the cloud, add to add files, and exit to exit the CLI.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Use != "init" && cmd.Use != "exit" && ctx == nil {
			return fmt.Errorf("database is not initialized. Please run the 'init' command first")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {

	RootCmd.Run(RootCmd, nil)
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		args := strings.Split(input, " ")
		if args[0] != "cloud" {
			fmt.Println("Please use 'cloud' prefix for commands.")
			continue
		}
		if len(args) < 2 {
			fmt.Println("Please use cloud [command]")
			continue
		}
		args = args[1:]
		RootCmd.SetArgs(args)

		if err := RootCmd.Execute(); err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println()
	}

}

func init() {
	RootCmd.PersistentFlags().BoolP("help", "h", false, "Help")
}
