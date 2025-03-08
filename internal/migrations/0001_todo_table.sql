-- +goose Up
CREATE TABLE todo
(
    id          BIGSERIAL PRIMARY KEY,
    name        text NOT NULL,
    description text,
    done        boolean
);


-- +goose Down
DROP TABLE todo;