-- +goose Up
-- +goose StatementBegin
CREATE TABLE todo_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    status INTEGER DEFAULT 0,
    createdat datetime default current_timestamp
);

INSERT INTO todo_new (id, name, status, createdat)
SELECT id, name, CASE WHEN done THEN 2 ELSE 0 END, createdat
FROM todo;

DROP TABLE todo;

ALTER TABLE todo_new RENAME TO todo;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE todo_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    done BOOLEAN DEFAULT false,
    createdat datetime default current_timestamp
);

INSERT INTO todo_old (id, name, done, createdat)
SELECT id, name, CASE WHEN status = 2 THEN 1 ELSE 0 END, createdat
FROM todo;

DROP TABLE todo;

ALTER TABLE todo_old RENAME TO todo;
-- +goose StatementEnd
