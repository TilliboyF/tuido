package db

import (
	"os"
	"testing"

	"github.com/TilliboyF/tuido/types"
	_ "github.com/mattn/go-sqlite3"
)

const (
	testDB = "./test.db"
)

func TestNewSqliteTodoStore(t *testing.T) {
	db, err := NewSqliteTodoStore()
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
	db, _ := NewSqliteTodoStore()
	todo := types.Todo{
		ID:   -1,
		Name: "test",
	}
	err := db.Add(&todo)
	if err != nil {
		t.Log("Error adding todo: ", err)
		t.Fail()
	}
	if todo.ID == -1 {
		t.Log("Id didn't get updated")
		t.Fail()
	}
	os.Remove(testDB)
}
