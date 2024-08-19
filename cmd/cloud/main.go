package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/cmd"
	"github.com/JIIL07/jcloud/internal/client/config"
	"github.com/JIIL07/jcloud/internal/client/lib/ctx"
	"github.com/JIIL07/jcloud/internal/client/models"
	cloud "github.com/JIIL07/jcloud/internal/client/operations"
	"github.com/JIIL07/jcloud/internal/client/storage"
	"log"
	"os"
	"strings"
)

func main() {
	c := config.MustLoad()

	// Initialize database
	s, err := storage.InitDatabase(c)
	if err != nil {
		log.Fatal(err)
	}

	fctx := &cloud.FileContext{
		Info:    &models.File{},
		Storage: &s,
	}

	ctx := context.Background()
	ctx = jctx.WithContext(ctx, "filecontext", fctx)
	cmd.SetContext(ctx)

	//// Execute the root command
	//if err := cmd.RootCmd.Execute(); err != nil {
	//	log.Fatal(err)
	//}

	reader := bufio.NewReader(os.Stdin)
	for {
		dir, _ := os.Getwd()
		fmt.Printf("%v>", dir)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		args := strings.Split(input, " ")
		cmd.RootCmd.SetArgs(args[1:])

		cmd.RootCmd.Execute()
	}
}
