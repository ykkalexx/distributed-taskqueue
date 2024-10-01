package task

import (
	"sort"
	"sync"
)

type Queue interface {
	AddTask(task Task) error
	GetTask() (Task, bool, error)
	Close() error
}

type InMemoryQueue struct {
	tasks []Task
	mu    sync.Mutex
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		tasks: make([]Task, 0),
	}
}

func (q *InMemoryQueue) AddTask(task Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tasks = append(q.tasks, task)
	sort.Slice(q.tasks, func(i, j int) bool {
		return q.tasks[i].Priority > q.tasks[j].Priority
	})
	return nil
}

func (q *InMemoryQueue) GetTask() (Task, bool, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.tasks) == 0 {
		return Task{}, false, nil
	}
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return task, true, nil
}

func (q *InMemoryQueue) Close() error {
	// No-op for in-memory queue
	return nil
}

// Ensure InMemoryQueue implements Queue interface
var _ Queue = (*InMemoryQueue)(nil)