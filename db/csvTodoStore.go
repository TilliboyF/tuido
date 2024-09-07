package db

import (
	"encoding/csv"
	"errors"
	"os"

	"github.com/TilliboyF/tuido/types"
)

type CsvTodoStore struct {
	file *os.File
}

func (s *CsvTodoStore) Close() {
	s.file.Close()
}

func NewCsvTodoStore(filePath string) (*CsvTodoStore, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			file, err = os.Create(filePath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &CsvTodoStore{
		file: file,
	}, nil
}

func (s *CsvTodoStore) determineNewIndex() (int64, error) {
	todos, err := s.GetAll()
	if err != nil {
		return 0, err
	}
	if len(todos) == 0 {
		return 0, nil
	}
	lasttodo := todos[len(todos)-1]
	return lasttodo.ID + 1, nil
}

func (s *CsvTodoStore) GetAll() ([]types.Todo, error) {
	s.file.Seek(0, 0)
	reader := csv.NewReader(s.file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := []types.Todo{}
	for _, row := range records {
		todo, err := types.NewTodoFromFields(row)
		if err != nil {
			return nil, err
		}
		result = append(result, todo)
	}
	return result, nil
}

func (s *CsvTodoStore) Add(todo types.Todo) (types.Todo, error) {

	newIndex, err := s.determineNewIndex()
	if err != nil {
		return types.Todo{}, err
	}
	todo.ID = newIndex

	writer := csv.NewWriter(s.file)
	defer writer.Flush()
	defer s.file.Sync()

	writer.Write(todo.CsvString())

	return todo, nil

}
