package main

import (
	"log"
	"time"

	"github.com/ykkalexx/distributed-taskqueue/internal/task"
	"github.com/ykkalexx/distributed-taskqueue/internal/worker"
	"github.com/ykkalexx/distributed-taskqueue/pkg/grpc"
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

	// Start gRPC server
	go func() {
		if err := grpc.StartServer(queue, 50051); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// give the second a moment to breath and start
	time.Sleep(time.Second)

	// submit task using grpc client
	functionNames := []string{"printHello", "simulateWork"}
	for i := 0; i < 10; i++ {
		err := grpc.SubmitTask("localhost:50051", int32(i), functionNames[i%len(functionNames)])
		if err != nil {
			log.Printf("Failed to submit task: %v", err)
		}
	}

	// Wait for tasks to complete
	time.Sleep(time.Second * 15)
}