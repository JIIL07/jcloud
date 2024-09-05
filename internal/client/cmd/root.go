package cmd

import (
	"context"
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/app"
	jctx "github.com/JIIL07/jcloud/pkg/ctx"
	"github.com/spf13/cobra"
	"os"
)

var (
	ctx         context.Context
	appCtx      *app.ClientContext
	versionFlag bool
	allFiles    bool
)

func SetContext(newCtx context.Context) {
	ctx = newCtx
}

var RootCmd = &cobra.Command{
	Use:     `jcloud`,
	Aliases: []string{"jc"},
	Short:   `Cloud is a cloud file manager CLI`,
	GroupID: "",
	Long: `Cloud is a cloud file manager CLI that provides various commands to manage files in the cloud.
It supports commands like init to initialize the cloud, add to add files, and exit to exit the CLI.`,
	Version: "0.0.1",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var ok bool
		appCtx, ok = jctx.FromContext[*app.ClientContext](ctx, "app-context")
		if !ok {
			return fmt.Errorf("failed to get file context")
		}

		//if cmd.Name() != loginCmd.Name() && cmd != cmd.Root() {
		//	content, err := io.ReadAll(appCtx.PathsService.P.JcloudFile)
		//	if err != nil {
		//		cobra.CheckErr(err)
		//	}
		//
		//	if len(content) == 0 {
		//		cobra.WriteStringAndCheck(os.Stdout, "use 'jcloud login [username] [email] [password]' first\n")
		//		os.Exit(0)
		//	}
		//
		//}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			cobra.WriteStringAndCheck(os.Stdout, cmd.Version)
		} else {
			err := cmd.Help()
			if err != nil {
				cobra.CheckErr(err)
			}
			return
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolP("help", "h", false, "Help")
	RootCmd.PersistentFlags().BoolVarP(&allFiles, "all", "a", false, "All files from the specified directory (or current directory if not specified)")
	RootCmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "Version")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		cobra.CheckErr(err)
	}
}
