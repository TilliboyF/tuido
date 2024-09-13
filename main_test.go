package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/TilliboyF/tuido/common"
	"github.com/TilliboyF/tuido/handler"
	"github.com/TilliboyF/tuido/types"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	handler, err := handler.NewTodoHandler(true)
	assert.NoError(t, err)
	rootCmd := NewRootCmd(handler)

	tests := []struct {
		name     string
		command  func() []string
		expected func() string
	}{
		{
			name: "Version print",
			command: func() []string {
				return []string{"version"}
			},
			expected: func() string {
				return fmt.Sprintf("tuido app version %s\n", version)
			},
		},
		{
			name: "Add Todo",
			command: func() []string {
				return []string{"add", "Task 1"}
			},
			expected: func() string {
				return common.TableStringFromTodo(types.Todo{
					Name:      "Task 1",
					ID:        1,
					Status:    types.TODO,
					CreatedAt: time.Now(),
				})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, err := common.ExecuteCommand(rootCmd, test.command()...)
			assert.NoError(t, err)
			assert.Equal(t, test.expected(), out)
		})
	}

}
