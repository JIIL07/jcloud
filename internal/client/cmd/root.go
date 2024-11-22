package cmd

import (
	"context"
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/app"
	"github.com/JIIL07/jcloud/internal/client/config"
	h "github.com/JIIL07/jcloud/internal/client/hints"
	"github.com/JIIL07/jcloud/pkg/ctx"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var (
	ctx         context.Context
	a           *app.ClientContext
	c           *config.ClientConfig
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
It supports commands like init to initialize the cloud, addFile to addFile files, and exit to exit the CLI.`,
	Version: "0.0.1",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var ok bool
		a, ok = jctx.FromContext[*app.ClientContext](ctx, "app-context")
		if !ok {
			return fmt.Errorf("failed to get file context")
		}

		c = a.Cfg

		if versionFlag || cmd.Name() == "login" || cmd == cmd.Root() {
			return nil
		}

		content, err := io.ReadAll(a.Paths.P.JcloudFile)
		if err != nil {
			return fmt.Errorf("failed to read login file: %v", err)
		}

		if len(content) == 0 {
			hintMessage := h.DisplayHint("login", h.LoginRequired, c)
			if hintMessage != "" {
				cobra.WriteStringAndCheck(os.Stdout, hintMessage)
			}
			os.Exit(0)
		}

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
	RootCmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "Version")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		cobra.CheckErr(err)
	}
}
