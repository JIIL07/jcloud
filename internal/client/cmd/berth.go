package cmd

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/jc"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/internal/client/requests"
	"github.com/JIIL07/jcloud/pkg/log"
	"github.com/spf13/cobra"
	"io"
)

var berthCmd = &cobra.Command{
	Use:   "berth",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		files := make([]models.File, 0)
		err := a.Storage.S.GetAllFiles(&files)
		if err != nil {
			cobra.CheckErr(err)
		}

		resp, err := requests.UploadFile(a, &files)
		if err != nil {
			cobra.CheckErr(err)
		}
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			cobra.CheckErr(err)
		}
		a.Logger.L.Info(fmt.Sprintf("response: %s, status code: %d", string(content), resp.StatusCode))
		cobra.CheckErr(resp.Body.Close())

		err = jc.DeleteAllFiles(a.File)
		if err != nil {
			a.Logger.L.Error("error deleting all files", jlog.Err(err))
			cobra.CheckErr(err)
		}

		a.Logger.L.Info("berth success")
	},
}

func init() {
	RootCmd.AddCommand(berthCmd)
}
