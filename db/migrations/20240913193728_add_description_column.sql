-- +goose Up
-- +goose StatementBegin
CREATE TABLE todo_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    description TEXT DEFAULT '',
    status INTEGER DEFAULT 0,
    createdat datetime default current_timestamp
);

INSERT INTO todo_new (id, name, description,status, createdat)
SELECT id, name, null,status, createdat
FROM todo;

DROP TABLE todo;

ALTER TABLE todo_new RENAME TO todo;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE todo_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    status INTEGER DEFAULT 0,
    createdat datetime default current_timestamp
);

INSERT INTO todo_old (id, name, status, createdat)
SELECT id, name, status, createdat
FROM todo;

DROP TABLE todo;

ALTER TABLE todo_old RENAME TO todo;
-- +goose StatementEnd
