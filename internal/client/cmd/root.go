package cmd

import (
	"context"
	"fmt"
	jctx "github.com/JIIL07/jcloud/internal/client/lib/ctx"
	cloud "github.com/JIIL07/jcloud/internal/client/operations"
	"github.com/spf13/cobra"
	"os"
)

var (
	fctx        *cloud.Context
	ctx         context.Context
	versionFlag bool
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
		fctx, ok = jctx.FromContext[*cloud.Context](ctx, "filecontext")
		if !ok {
			return fmt.Errorf("failed to get file context")
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			cobra.WriteStringAndCheck(os.Stdout, cmd.Version)
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolP("help", "h", false, "Help")
	RootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Version")
}
