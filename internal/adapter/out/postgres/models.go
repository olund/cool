// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package postgres

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Author struct {
	ID   int64
	Name string
	Bio  pgtype.Text
}
