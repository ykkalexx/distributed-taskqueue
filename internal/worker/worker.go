package worker

import (
	"fmt"
	"time"

	"github.com/ykkalexx/distributed-taskqueue/internal/logger"
	"github.com/ykkalexx/distributed-taskqueue/internal/task"
)

type Worker struct {
	id    int
	queue task.Queue
}

func NewWorker(id int, queue task.Queue) *Worker {
	return &Worker{
		id:    id,
		queue: queue,
	}
}

func (w *Worker) ID() int {
	return w.id
}

func (w *Worker) Start() {
	logger.Info("Worker %d started", w.id)
	for {
		t, ok, err := w.queue.GetTask()
		if err != nil {
			logger.Error("Worker %d error getting task: %v", w.id, err)
			time.Sleep(time.Second)
			continue
		}
		if !ok {
			time.Sleep(time.Second)
			continue
		}

		logger.Info("Worker %d executing task %d (%s) - Attempt %d", w.id, t.ID, t.FunctionName, t.Retries+1)

		if fn, exists := FunctionMap[t.FunctionName]; exists {
			err = fn()
			if err != nil {
				logger.Warn("Error executing task %d: %v", t.ID, err)
				if t.Retries < t.MaxRetries {
					t.Retries++
					logger.Info("Retrying task %d (Attempt %d/%d)", t.ID, t.Retries+1, t.MaxRetries+1)
					err = w.queue.AddTask(t)
					if err != nil {
						logger.Error("Failed to requeue task %d: %v", t.ID, err)
					}
				} else {
					logger.Warn("Task %d failed after %d attempts", t.ID, t.MaxRetries+1)
				}
			} else {
				logger.Info("Task %d completed successfully", t.ID)
			}
		} else {
			logger.Warn("Unknown function for task %d: %s", t.ID, t.FunctionName)
		}
	}
}

var FunctionMap = map[string]func() error{
	"printHello": func() error {
		fmt.Println("Hello")
		return nil
	},
	"simulateWork": func() error {
		time.Sleep(time.Second)
		return nil
	},
	// will add more functions
}
