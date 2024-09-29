package main

import (
	"fmt"
	"time"

	"github.com/ykkalexx/distributed-taskqueue/internal/task"
	"github.com/ykkalexx/distributed-taskqueue/internal/worker"
)

func main() {
	queue := task.NewQueue()

	// start some workers
	for i := 0; i < 3; i++ {
		go worker.Start(i, queue)
	}

	// add some tasks
	for i := 0; i < 10; i++ {
		taskID := i
		queue.AddTask(task.Task{
			ID: taskID,
			Function: func() error {
				fmt.Printf("Executing task %d\n", taskID)
				time.Sleep(time.Second) // Simulate work
				return nil
			},
		})
	}

	// Wait for tasks to complete
	time.Sleep(time.Second * 15)
}