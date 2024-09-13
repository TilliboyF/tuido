package db

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TilliboyF/tuido/types"
	"github.com/stretchr/testify/assert"
)

func TestSqliteTodoStore_InMemory(t *testing.T) {

	store, _, err := NewSqliteTodoStore(false, true)
	assert.NoError(t, err)
	defer store.Close()

	tests := []struct {
		name      string
		action    func() error
		validate  func(t *testing.T)
		expectErr bool
	}{
		{
			name: "Add todo item",
			action: func() error {
				todo := types.Todo{Name: "Task 1"}
				return store.Add(&todo)
			},
			validate: func(t *testing.T) {
				todos, err := store.GetAll()
				assert.NoError(t, err)
				assert.Equal(t, 1, len(todos))
				assert.Equal(t, "Task 1", todos[0].Name)
			},
		},
		{
			name: "Get all todo items",
			action: func() error {
				// Already added a todo in previous test, so we just need to validate
				return nil
			},
			validate: func(t *testing.T) {
				// Validate the todo was fetched correctly
				todos, err := store.GetAll()
				assert.NoError(t, err)
				assert.Equal(t, 1, len(todos))
				assert.Equal(t, "Task 1", todos[0].Name)
				assert.Equal(t, types.TODO, todos[0].Status)
			},
		},
		{
			name: "Get todo by ID",
			action: func() error {
				// Fetch the todo by its ID (ID 1 was set by previous tests)
				_, err := store.GetById(1)
				return err
			},
			validate: func(t *testing.T) {
				todo, err := store.GetById(1)
				assert.NoError(t, err)
				assert.Equal(t, "Task 1", todo.Name)
				assert.Equal(t, types.TODO, todo.Status)
			},
		},
		{
			name: "Complete todo item",
			action: func() error {
				// Complete the todo (set Done to true)
				return store.Complete(1)
			},
			validate: func(t *testing.T) {
				// Verify that the todo is marked as done
				todo, err := store.GetById(1)
				assert.NoError(t, err)
				assert.Equal(t, types.DONE, todo.Status)
			},
		},
		{
			name: "Delete todo item",
			action: func() error {
				// Delete the todo
				return store.Delete(1)
			},
			validate: func(t *testing.T) {
				// Verify that the todo was deleted
				_, err := store.GetById(1)
				assert.Error(t, err) // Expect an error since the item should be deleted
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.action()
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if test.validate != nil {
				test.validate(t)
			}
		})
	}
}

func TestSqliteTodoStore_Mock(t *testing.T) {
	store, mock, err := NewSqliteTodoStore(true, false)
	assert.NoError(t, err)
	defer store.Close()

	tests := []struct {
		name      string
		mock      func()
		validate  func(t *testing.T)
		expectErr bool
	}{
		{
			name: "Add todo item",
			mock: func() {
				mock.ExpectPrepare("INSERT INTO todo").
					ExpectExec().
					WithArgs("Task 1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T) {
				todo := &types.Todo{Name: "Task 1"}
				err := store.Add(todo)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), todo.ID)
			},
		},
		{
			name: "Get all todo items",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "done", "createdat"}).
					AddRow(1, "Task 1", 0, time.Now()).
					AddRow(2, "Task 2", 0, time.Now())
				mock.ExpectPrepare("SELECT \\* FROM todo").
					ExpectQuery().
					WillReturnRows(rows)
			},
			validate: func(t *testing.T) {
				todos, err := store.GetAll()
				assert.NoError(t, err)
				assert.Equal(t, 2, len(todos))
				assert.Equal(t, "Task 1", todos[0].Name)
				assert.Equal(t, "Task 2", todos[1].Name)
			},
		},
		{
			name: "Get todo by ID",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "done", "createdat"}).
					AddRow(1, "Task 1", 0, time.Now())
				mock.ExpectPrepare("SELECT \\* FROM todo WHERE id=\\?").
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(rows)
			},
			validate: func(t *testing.T) {
				todo, err := store.GetById(1)
				assert.NoError(t, err)
				assert.Equal(t, "Task 1", todo.Name)
			},
		},
		{
			name: "Complete todo item",
			mock: func() {
				// Mock the behavior for Complete (UPDATE query)
				mock.ExpectExec("UPDATE todo SET status=2 WHERE id=\\?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			validate: func(t *testing.T) {
				// Perform the action and validate
				err := store.Complete(1)
				assert.NoError(t, err)
			},
		},
		{
			name: "Delete todo item",
			mock: func() {
				// Mock the behavior for Delete (DELETE query)
				mock.ExpectExec("DELETE FROM todo WHERE id=\\?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			validate: func(t *testing.T) {
				// Perform the action and validate
				err := store.Delete(1)
				assert.NoError(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.mock != nil {
				test.mock()
			}
			test.validate(t)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
