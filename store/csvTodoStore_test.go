package store

import (
	"os"
	"testing"
)

const (
	TestFileName = "test.csv"
)

func TestNewCsvTodoStore(t *testing.T) {

	_, err := NewCsvTodoStore(TestFileName)
	if err != nil {
		t.Log("Couldn't create store: ", err)
		t.Fail()
	}

	cleanup()
}

func cleanup() {
	os.Remove(TestFileName)
}

func TestCsvStoreAdd(t *testing.T) {
	store, _ := NewCsvTodoStore(TestFileName)
	todo := Todo{
		ID:   -1,
		Name: "jkwdbwjkvbwv",
		Done: false,
	}
	updatedTodo, err := store.Add(todo)
	if err != nil {
		t.Log("Error while adding todo: ", err)
		t.Fail()
	}
	if updatedTodo.ID == -1 {
		t.Log("Add didn't update todo")
		t.Fail()
	}

	cleanup()

}
func TestCsvStoreGetAll(t *testing.T) {
	store, _ := NewCsvTodoStore(TestFileName)
	todo := Todo{
		ID:   -1,
		Name: "Test description",
		Done: false,
	}

	for i := 0; i < 10; i++ {
		store.Add(todo)
	}

	todos, err := store.GetAll()
	if err != nil {
		t.Log("Couldn't get todos: ", err)
	}

	if len(todos) != 10 {
		t.Log("todos don't are the right amount")
	}

	cleanup()
}
