package cmd

import (
	"fmt"
	slg "github.com/JIIL07/jcloud/internal/client/lib/logger"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var whoiamCmd = &cobra.Command{
	Use:     "whoiam",
	Aliases: []string{"wia"},
	Short:   "Return current user",
	Long:    "Read .jcloud file from $HOME dir and return information about current user",
	Run: func(cmd *cobra.Command, args []string) {
		var bytes []byte
		bytes, err := os.ReadFile(paths.Jcloud.Name())
		if err != nil {
			logger.Error("error reading file", slg.Err(err))
			cobra.CheckErr(err)
		}
		me := strings.Split(string(bytes), " ")
		fmt.Printf("username: %v\nemail: %v\n", me[0], me[1])
	},
}

func init() {
	RootCmd.AddCommand(whoiamCmd)
}
