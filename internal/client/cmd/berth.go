package cmd

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/internal/client/requests"
	"github.com/spf13/cobra"
	"io"
)

// berthCmd represents the berth command
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
		err := appCtx.StorageService.S.GetAllFiles(&files)
		if err != nil {
			cobra.CheckErr(err)
		}

		resp, err := requests.UploadFile(&files)
		if err != nil {
			cobra.CheckErr(err)
		}
		c, err := io.ReadAll(resp.Body)
		if err != nil {
			cobra.CheckErr(err)
		}
		appCtx.LoggerService.L.Info(fmt.Sprintf("response: %s, status code: %d", string(c), resp.StatusCode))
		cobra.CheckErr(resp.Body.Close())
	},
}

func init() {
	RootCmd.AddCommand(berthCmd)
}
