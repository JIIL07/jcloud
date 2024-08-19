package cmd

import (
	"context"
	"fmt"
	jctx "github.com/JIIL07/jcloud/internal/client/lib/ctx"
	cloud "github.com/JIIL07/jcloud/internal/client/operations"
	"github.com/spf13/cobra"
)

var (
	fctx *cloud.FileContext
	ctx  context.Context
)

func SetContext(newCtx context.Context) {
	ctx = newCtx
}

var RootCmd = &cobra.Command{
	Use:   `cloud`,
	Short: `Cloud is a cloud file manager CLI`,
	Long: `Cloud is a cloud file manager CLI that provides various commands to manage files in the cloud.
It supports commands like init to initialize the cloud, add to add files, and exit to exit the CLI.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var ok bool
		fctx, ok = jctx.FromContext[*cloud.FileContext](ctx, "filecontext")
		if !ok {
			return fmt.Errorf("failed to get file context")
		}

		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().BoolP("help", "h", false, "Help")
}
