// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package todo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createTodo = `-- name: CreateTodo :one
INSERT INTO todo (
    name, description
) VALUES (
             $1, $2
         )
    RETURNING id, name, description, done
`

type CreateTodoParams struct {
	Name        string
	Description pgtype.Text
}

func (q *Queries) CreateTodo(ctx context.Context, arg CreateTodoParams) (Todo, error) {
	row := q.db.QueryRow(ctx, createTodo, arg.Name, arg.Description)
	var i Todo
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Done,
	)
	return i, err
}

const deleteTodo = `-- name: DeleteTodo :exec
DELETE FROM todo
WHERE id = $1
`

func (q *Queries) DeleteTodo(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteTodo, id)
	return err
}

const getTodo = `-- name: GetTodo :one
SELECT id, name, description, done FROM todo
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTodo(ctx context.Context, id int64) (Todo, error) {
	row := q.db.QueryRow(ctx, getTodo, id)
	var i Todo
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Done,
	)
	return i, err
}

const listTodos = `-- name: ListTodos :many
SELECT id, name, description, done FROM todo
ORDER BY name
`

func (q *Queries) ListTodos(ctx context.Context) ([]Todo, error) {
	rows, err := q.db.Query(ctx, listTodos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Todo
	for rows.Next() {
		var i Todo
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Done,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTodo = `-- name: UpdateTodo :exec
UPDATE todo
set name = $2,
    description = $3,
    done = $4
WHERE id = $1
`

type UpdateTodoParams struct {
	ID          int64
	Name        string
	Description pgtype.Text
	Done        pgtype.Bool
}

func (q *Queries) UpdateTodo(ctx context.Context, arg UpdateTodoParams) error {
	_, err := q.db.Exec(ctx, updateTodo,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Done,
	)
	return err
}
