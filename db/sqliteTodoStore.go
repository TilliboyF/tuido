package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TilliboyF/tuido/common"
	"github.com/TilliboyF/tuido/types"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type SqliteTodoStore struct {
	db *sql.DB
}

func getDBPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dbPath := filepath.Join(configDir, "tuido", common.DB_NAME)
	return dbPath, nil
}

func initializeDB(useInMemory bool) (*sql.DB, error) {

	var dbPath string
	var err error

	if useInMemory {
		dbPath = ":memory:"
	} else {
		dbPath, err = getDBPath()
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(embedMigrations)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return nil, err
	}
	return db, nil
}

func NewSqliteTodoStore(useMock bool, useInMemory bool) (*SqliteTodoStore, sqlmock.Sqlmock, error) {

	var db *sql.DB
	var mock sqlmock.Sqlmock
	var err error

	if useMock {
		db, mock, err = sqlmock.New()
		if err != nil {
			return nil, nil, err
		}
	} else {
		db, err = initializeDB(useInMemory)
		if err != nil {
			return nil, nil, err
		}
	}

	data := SqliteTodoStore{
		db: db,
	}
	return &data, mock, nil
}

func (s *SqliteTodoStore) Add(t *types.Todo) error {
	query := `INSERT INTO todo (name) VALUES (?);`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(t.Name)
	if err != nil {
		return err
	}
	insertedID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = insertedID
	return nil
}

func (s *SqliteTodoStore) GetAll() ([]types.Todo, error) {
	query := `SELECT * FROM todo;`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []types.Todo

	for rows.Next() {
		var todo types.Todo
		if err := rows.Scan(&todo.ID, &todo.Name, &todo.Done, &todo.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (s *SqliteTodoStore) GetById(id int64) (types.Todo, error) {
	query := `SELECT * FROM todo WHERE id=?`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return types.Todo{}, err
	}
	defer stmt.Close()
	var todo types.Todo

	row := stmt.QueryRow(id)
	if err := row.Scan(&todo.ID, &todo.Name, &todo.Done, &todo.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.Todo{}, fmt.Errorf("no todo found with id=%d", id)
		}
		return types.Todo{}, err
	}

	return todo, nil
}

func (s *SqliteTodoStore) Complete(id int64) error {
	query := "UPDATE todo SET done=true WHERE id=?;"
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqliteTodoStore) Delete(id int64) error {
	query := `DELETE FROM todo WHERE id=?;`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *SqliteTodoStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
