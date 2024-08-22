package cmd

import (
	"fmt"
	slg "github.com/JIIL07/jcloud/internal/client/lib/logger"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// whoiamCmd represents the whoiam command
var whoiamCmd = &cobra.Command{
	Use:     "whoiam",
	Aliases: []string{"wia"},
	Short:   "Return current user",
	Long:    "Read .jcloud file from $HOME dir and return information about current user",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(fctx.Local.Jcloud)
		if err != nil {
			fctx.Logger.Error("error opening file", slg.Err(err))
			cobra.CheckErr(err)
		}
		defer f.Close()

		var bytes []byte
		bytes, err = os.ReadFile(f.Name())
		if err != nil {
			fctx.Logger.Error("error reading file", slg.Err(err))
			cobra.CheckErr(err)
		}
		me := strings.Split(string(bytes), " ")
		fmt.Printf("username: %v\nemail: %v\n", me[0], me[1])
	},
}

func init() {
	RootCmd.AddCommand(whoiamCmd)
}
