package store

import (
	"os"
	"testing"
)

const (
	testDB = "./test.db"
)

func TestNewSqliteTodoStore(t *testing.T) {
	db, err := NewSqliteTodoStore(testDB)
	if nil != err {
		t.Log("Error creating db: ", err)
		t.Fail()
	}
	if db.db == nil {
		t.Log("db not initialised!")
		t.Fail()
	}
	os.Remove(testDB)
}

func TestSqliteStoreAdd(t *testing.T) {
	db, _ := NewSqliteTodoStore(testDB)
	updatedTodo, err := db.Add(Todo{
		ID:   -1,
		Name: "test",
	})
	if err != nil {
		t.Log("Error adding todo: ", err)
		t.Fail()
	}
	if updatedTodo.ID == -1 {
		t.Log("Id didn't get updated")
		t.Fail()
	}
	os.Remove(testDB)
}
