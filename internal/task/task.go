package task

type Priority int

const (
	LowPriority Priority = iota
	MediumPriority
	HighPriority
)

type Task struct {
	ID           int          `json:"id"`
	FunctionName string       `json:"function_name"`
	Priority 	 Priority     `json:"priority"`
}