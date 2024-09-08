package main

import (
	"fmt"
	"testing"

	"github.com/TilliboyF/tuido/common"
	"github.com/TilliboyF/tuido/handler"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	handler, err := handler.NewTodoHandler(true)
	assert.NoError(t, err)
	rootCmd := NewRootCmd(handler)

	out, err := common.ExecuteCommand(rootCmd, "version")
	assert.Equal(t, fmt.Sprintf("tuido app version %s\n", version), out)
}
