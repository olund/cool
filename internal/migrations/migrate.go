package migrations

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/jackc/pgx/v5/stdlib"
)

func Run(ctx context.Context, postgresConnectionString, migrationDir string) error {

	dbpool, err := pgxpool.New(context.Background(), postgresConnectionString)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	db := stdlib.OpenDBFromPool(dbpool)

	if err := goose.Version(db, migrationDir); err != nil {
		return fmt.Errorf("goose version before: %w", err)
	}

	if err := goose.Up(db, migrationDir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	if err := goose.Version(db, migrationDir); err != nil {
		return fmt.Errorf("goose version after: %w", err)
	}

	return nil
}
