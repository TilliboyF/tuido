package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteTodoStore struct {
	db *sql.DB
}

func NewSqliteTodoStore(connectString string) (*SqliteTodoStore, error) {
	db, err := sql.Open("sqlite3", connectString)
	if err != nil {
		return nil, err
	}
	data := SqliteTodoStore{
		db: db,
	}
	// check if seeding is needed
	_, err = db.Query("SELECT * FROM todo")
	if err != nil {
		if err := data.Seed(); err != nil {
			log.Println(err)
		}
	}
	return &data, nil

}

func (s *SqliteTodoStore) Add(t Todo) (Todo, error) {

	query := `INSERT INTO todo (name) VALUES (?);`
	stmt, err := s.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return Todo{}, err
	}
	result, err := stmt.Exec(t.Name)
	if err != nil {
		return Todo{}, err
	}
	insertedID, err := result.LastInsertId()
	if err != nil {
		return Todo{}, err
	}
	t.ID = insertedID
	return t, nil
}

func (s *SqliteTodoStore) GetAll() ([]Todo, error) {
	query := `Select * from todo;`
	stmt, err := s.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Name, &todo.Done, &todo.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (s *SqliteTodoStore) GetById(id int64) (Todo, error) {
	query := `Select * FROM todo WHERE id=?`
	stmt, err := s.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return Todo{}, err
	}

	var todo Todo

	row := stmt.QueryRow(id)
	if err := row.Scan(&todo.ID, &todo.Name, &todo.Done, &todo.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, fmt.Errorf("No todo found with id=%d", id)
		}
		return Todo{}, err
	}

	return todo, nil
}

func (s *SqliteTodoStore) Seed() error {
	slog.Info("Seeding db...")

	stmt := `CREATE TABLE todo (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT,
				done BOOLEAN DEFAULT false,
				createdat datetime default current_timestamp
			)`

	_, err := s.db.Exec(stmt)
	return err
}
