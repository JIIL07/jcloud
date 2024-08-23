package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/cmd"
	"github.com/JIIL07/jcloud/internal/client/config"
	"github.com/JIIL07/jcloud/internal/client/lib/ctx"
	"github.com/JIIL07/jcloud/internal/client/lib/home"
	"github.com/JIIL07/jcloud/internal/client/models"
	cloud "github.com/JIIL07/jcloud/internal/client/operations"
	"github.com/JIIL07/jcloud/internal/client/storage"
	"log"
	"log/slog"
	"os"
	"strings"
)

func main() {
	c := config.MustLoad()

	s, err := storage.InitDatabase(c)
	if err != nil {
		log.Fatal(err)
	}

	fctx := &cloud.Context{
		File:    &models.File{},
		Storage: &s,
	}

	paths := home.SetPaths()
	defer paths.Close()

	logger := slog.New(slog.NewTextHandler(paths.Jlog, &slog.HandlerOptions{Level: slog.LevelDebug}))

	ctx := jctx.WithContext(context.Background(), "context", fctx)
	ctx = jctx.WithContext(ctx, "logger", logger)
	ctx = jctx.WithContext(ctx, "paths", paths)
	cmd.SetContext(ctx)

	switch {
	case c.Env == "local":
		if err := cmd.RootCmd.Execute(); err != nil {
			log.Fatal(err)
		}
	case c.Env == "debug":
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
}
