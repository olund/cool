package migrations

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Run(db *sql.DB, migrationDir string) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

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
