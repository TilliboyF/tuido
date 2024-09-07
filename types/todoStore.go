package types

type TodoStore interface {
	Add(*Todo) error
	GetById(int64) (Todo, error)
	GetAll() ([]Todo, error)
	Delete(int64) error
	Complete(int64) error
}
