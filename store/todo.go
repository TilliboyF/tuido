package store

import (
	"fmt"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mergestat/timediff"
)

type Todo struct {
	ID        int64
	Name      string
	Done      bool
	CreatedAt time.Time
}

func (t *Todo) String() string {
	return fmt.Sprintf("%d|%s|%v|%v", t.ID, t.Name, "false", timediff.TimeDiff(t.CreatedAt))
}

func NewTodoFromFields(fields []string) (Todo, error) {
	id, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return Todo{}, err
	}
	var done bool
	if fields[2] == "1" {
		done = true
	} else {
		done = false
	}

	return Todo{
		ID:   id,
		Name: fields[1],
		Done: done,
	}, nil
}

func (t *Todo) CsvString() []string {
	var done string
	if t.Done {
		done = "1"
	} else {
		done = "0"
	}
	return []string{fmt.Sprint(t.ID), t.Name, done}
}

type TodoStore interface {
	Add(Todo) (Todo, error)
	Update(Todo) (Todo, error)
	GetById(int64) (Todo, error)
	GetAll() ([]Todo, error)
	Delete(Todo) error
}
