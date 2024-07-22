<<<<<<< e11e69cc92d7c5a0aef5e57c3c44bea1b6154e12:cmd/whoiam.go
/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// whoiamCmd represents the whoiam command
var whoiamCmd = &cobra.Command{
	Use:   "whoiam",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		loadUserPass()
	},
}

func init() {
	RootCmd.AddCommand(whoiamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whoiamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whoiamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
=======
package cmd

import (
	"github.com/spf13/cobra"
)

// whoiamCmd represents the whoiam command
var whoiamCmd = &cobra.Command{
	Use:   "whoiam",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		loadUserPass()
	},
}

func init() {
	RootCmd.AddCommand(whoiamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whoiamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whoiamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
>>>>>>> Big file structure update, dockerfile does not work currently, some new features in code: new server system (not complete) and new logger:internal/cmd/whoiam.go
