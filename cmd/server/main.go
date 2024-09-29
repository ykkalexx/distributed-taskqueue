package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ykkalexx/distributed-taskqueue/internal/task"
	"github.com/ykkalexx/distributed-taskqueue/internal/worker"
)

func main() {
	// Create a Redis-based queue
	queue, err := task.NewRedisQueue("localhost:6379", "", 0, "tasks")
	if err != nil {
		log.Fatalf("Failed to create Redis queue: %v", err)
	}
	defer queue.Close()

	// start some workers
	for i := 0; i < 3; i++ {
		go worker.Start(i, queue)
	}

	// add some tasks
	for i := 0; i < 10; i++ {
		taskID := i
		err := queue.AddTask(task.Task{
			ID: taskID,
			Function: func() error {
				fmt.Printf("Executing task %d\n", taskID)
				time.Sleep(time.Second) // Simulate work
				return nil
			},
		})
		if err != nil {
			log.Printf("Failed to add task: %v", err)
		}
	}

	// Wait for tasks to complete
	time.Sleep(time.Second * 15)
}