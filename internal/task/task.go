package task

type Task struct {
	ID       int
	Function func() error
}