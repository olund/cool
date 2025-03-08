-- name: GetTodo :one
SELECT * FROM todo
WHERE id = $1 LIMIT 1;

-- name: ListTodos :many
SELECT * FROM todo
ORDER BY name;

-- name: CreateTodo :one
INSERT INTO todo (
    name, description
) VALUES (
             $1, $2
         )
    RETURNING *;

-- name: UpdateTodo :exec
UPDATE todo
set name = $2,
    description = $3,
    done = $4
WHERE id = $1;

-- name: DeleteTodo :exec
DELETE FROM todo
WHERE id = $1;