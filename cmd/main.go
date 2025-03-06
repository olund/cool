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
		case "POSTGRES_CONNECTION_STRING":
			return "postgresql://postgres:mysecretpassword@localhost:5432/postgres"
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
