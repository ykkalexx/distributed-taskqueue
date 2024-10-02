package worker

import (
	"fmt"
	"log"
	"time"

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
	for {
		t, ok, err := w.queue.GetTask()
		if err != nil {
			log.Printf("Worker %d error getting task: %v", w.id, err)
			time.Sleep(time.Second)
			continue
		}
		if !ok {
			time.Sleep(time.Second)
			continue
		}

		fmt.Printf("Worker %d executing task %d (%s)\n", w.id, t.ID, t.FunctionName)

		if fn, exists := FunctionMap[t.FunctionName]; exists {
			err = fn()
			if err != nil {
				fmt.Printf("Error executing task %d: %v\n", t.ID, err)
				if t.Retries < t.MaxRetries {
					t.Retries++
					fmt.Printf("Retrying task %d (Attempt %d/%d)\n", t.ID, t.Retries+1, t.MaxRetries+1)
					err = w.queue.AddTask(t)
					if err != nil {
						fmt.Printf("Failed to requeue task %d: %v\n", t.ID, err)
					}
				} else {
					fmt.Printf("Task %d failed after %d attempts\n", t.ID, t.MaxRetries+1)
				}
			} else {
				fmt.Printf("Task %d completed successfully\n", t.ID)
			}
		} else {
			fmt.Printf("Unknown function for task %d: %s\n", t.ID, t.FunctionName)
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
