package main

import (
	"log"
	"time"

	"github.com/ykkalexx/distributed-taskqueue/internal/loadbalancer"
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

	// Create load balancer
	lb := loadbalancer.New()

	// Create and start workers
	for i := 0; i < 3; i++ {
		w := worker.NewWorker(i, queue)
		lb.AddWorker(w)
		go w.Start()
	}

	// Start gRPC server with load balancer
	go func() {
		if err := grpc.StartServer(queue, lb, 50051); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// give the server time to wake up and get coffee zzzZzzz
	time.Sleep(time.Second)

	// Submit some tasks using the gRPC client
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