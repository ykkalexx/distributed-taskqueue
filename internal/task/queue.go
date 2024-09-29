package task

import "sync"

type Queue struct {
	tasks []Task
	mu    sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		tasks: make([]Task, 0),
	}
}

func (q *Queue) AddTask(task Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tasks = append(q.tasks, task)
}

func (q *Queue) GetTask() (Task, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.tasks) == 0 {
		return Task{}, false
	}
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return task, true
}