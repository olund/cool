-- name: GetTodo :one
SELECT * FROM todo
WHERE id = ? LIMIT 1;

-- name: ListTodos :many
SELECT * FROM todo
ORDER BY name;

-- name: CreateTodo :one
INSERT INTO todo (
    name, description
) VALUES (
             ?, ?
         )
    RETURNING *;

-- name: UpdateTodo :exec
UPDATE todo
set name = ?,
    description = ?,
    done = ?
WHERE id = ?;

-- name: DeleteTodo :exec
DELETE FROM todo
WHERE id = ?;