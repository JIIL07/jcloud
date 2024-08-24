package cmd

import (
	"context"
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/app"
	jctx "github.com/JIIL07/jcloud/internal/client/lib/ctx"
	"github.com/JIIL07/jcloud/internal/client/lib/home"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

var (
	ctx         context.Context
	fs          *app.ClientContext
	paths       *home.Paths
	logger      *slog.Logger
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
		fs, ok = jctx.FromContext[*app.ClientContext](ctx, "app-context")
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
