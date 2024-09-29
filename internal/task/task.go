package task

type Task struct {
	ID       int          `json:"id"`
	Function func() error `json:"-"`
}