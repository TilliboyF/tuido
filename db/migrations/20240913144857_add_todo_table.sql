-- +goose Up
-- +goose StatementBegin
CREATE TABLE todo (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			done BOOLEAN DEFAULT false,
			createdat datetime default current_timestamp
		);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE todo;
-- +goose StatementEnd
