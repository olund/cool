package main

import (
	"context"
	"fmt"
	"github.com/olund/cool/internal"
	"github.com/olund/cool/internal/config"
	"os"
)

func main() {
	ctx := context.Background()

	app := internal.NewApp()

	cfg := config.Config{
		Host: "localhost",
		Port: "8080",
	}

	getenv := func(s string) string {
		switch s {
		case "MIGRATIONS_DIR":
			return "../internal/migrations"
		case "DB_NAME":
			return "todo/todo.db"

		}

		return s
	}

	if err := app.Run(ctx, os.Stdout, getenv, cfg); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "%s\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
