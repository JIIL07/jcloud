// nolint:errcheck
package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/app"
	"github.com/JIIL07/jcloud/internal/client/cmd"
	"github.com/JIIL07/jcloud/internal/client/config"
	"github.com/JIIL07/jcloud/internal/client/lib/ctx"
	"log"
	"os"
	"strings"
)

func main() {
	c := config.MustLoad()

	appc, err := app.NewAppContext(c)
	if err != nil {
		log.Fatal(err)
	}

	ctx := jctx.WithContext(context.Background(), "app-context", appc)
	cmd.SetContext(ctx)

	switch {
	case c.Env == "local":
		cmd.Execute()
	case c.Env == "debug":
		reader := bufio.NewReader(os.Stdin)
		for {
			dir, _ := os.Getwd()
			fmt.Printf("%v>", dir)

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			args := strings.Split(input, " ")
			cmd.RootCmd.SetArgs(args[1:])

			cmd.Execute()
		}
	}
}
