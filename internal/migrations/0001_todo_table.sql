-- +goose Up
CREATE TABLE todo
(
    id          INTEGER PRIMARY KEY,
    name        text NOT NULL,
    description text,
    done        boolean
);


-- +goose Down
DROP TABLE todo;