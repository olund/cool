package migrations

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Run(ctx context.Context, db *sql.DB, migrationDir string) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.UpContext(ctx, db, migrationDir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
