package types

import (
	"fmt"
	"time"

	"github.com/mergestat/timediff"
)

type Status int

const (
	TODO Status = iota
	INPROGRESS
	DONE
)

func (s Status) String() string {
	switch s {
	case TODO:
		return "todo"
	case INPROGRESS:
		return "in progress"
	case DONE:
		return "done"
	}
	return ""
}

type Todo struct {
	ID          int64
	Name        string
	Description string
	Status      Status
	CreatedAt   time.Time
}

func (t Todo) String() string {
	return fmt.Sprintf("%d|%s|%s|%v", t.ID, t.Name, t.Status, timediff.TimeDiff(t.CreatedAt))
}
