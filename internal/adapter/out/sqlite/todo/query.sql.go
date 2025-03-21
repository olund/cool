// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package todo

import (
	"context"
	"database/sql"
)

const createTodo = `-- name: CreateTodo :one
INSERT INTO todo (
    name, description
) VALUES (
             ?, ?
         )
    RETURNING id, name, description, done
`

type CreateTodoParams struct {
	Name        string
	Description sql.NullString
}

func (q *Queries) CreateTodo(ctx context.Context, arg CreateTodoParams) (Todo, error) {
	row := q.db.QueryRowContext(ctx, createTodo, arg.Name, arg.Description)
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
WHERE id = ?
`

func (q *Queries) DeleteTodo(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTodo, id)
	return err
}

const getTodo = `-- name: GetTodo :one
SELECT id, name, description, done FROM todo
WHERE id = ? LIMIT 1
`

func (q *Queries) GetTodo(ctx context.Context, id int64) (Todo, error) {
	row := q.db.QueryRowContext(ctx, getTodo, id)
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
ORDER BY id DESC
`

func (q *Queries) ListTodos(ctx context.Context) ([]Todo, error) {
	rows, err := q.db.QueryContext(ctx, listTodos)
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
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTodo = `-- name: UpdateTodo :exec
UPDATE todo
set name = ?,
    description = ?,
    done = ?
WHERE id = ?
`

type UpdateTodoParams struct {
	Name        string
	Description sql.NullString
	Done        sql.NullBool
	ID          int64
}

func (q *Queries) UpdateTodo(ctx context.Context, arg UpdateTodoParams) error {
	_, err := q.db.ExecContext(ctx, updateTodo,
		arg.Name,
		arg.Description,
		arg.Done,
		arg.ID,
	)
	return err
}

const updateTodoDoneState = `-- name: UpdateTodoDoneState :exec
UPDATE todo
set done = ?
WHERE id = ?
`

type UpdateTodoDoneStateParams struct {
	Done sql.NullBool
	ID   int64
}

func (q *Queries) UpdateTodoDoneState(ctx context.Context, arg UpdateTodoDoneStateParams) error {
	_, err := q.db.ExecContext(ctx, updateTodoDoneState, arg.Done, arg.ID)
	return err
}
