package main

import (
	"log"
	"time"

	"github.com/ykkalexx/distributed-taskqueue/internal/loadbalancer"
	"github.com/ykkalexx/distributed-taskqueue/internal/task"
	"github.com/ykkalexx/distributed-taskqueue/internal/worker"
	"github.com/ykkalexx/distributed-taskqueue/pkg/grpc"
)

type CompositeQueue struct {
	redisQueue  *task.RedisQueue
	sqliteQueue *task.SQLiteQueue
}

func (cq *CompositeQueue) AddTask(t task.Task) error {
	if err := cq.redisQueue.AddTask(t); err != nil {
		return err
	}
	return cq.sqliteQueue.AddTask(t)
}

func (cq *CompositeQueue) GetTask() (task.Task, bool, error) {
	t, ok, err := cq.redisQueue.GetTask()
	if err != nil || ok {
		return t, ok, err
	}
	return cq.sqliteQueue.GetTask()
}

func (cq *CompositeQueue) Close() error {
	if err := cq.redisQueue.Close(); err != nil {
		return err
	}
	return cq.sqliteQueue.Close()
}

func main() {
	redisQueue, err := task.NewRedisQueue("localhost:6379", "", 0, "tasks")
	if err != nil {
		log.Fatalf("Failed to create Redis queue: %v", err)
	}

	sqliteQueue, err := task.NewSQLiteQueue("tasks.db")
	if err != nil {
		log.Fatalf("Failed to create SQLite queue: %v", err)
	}

	queue := &CompositeQueue{
		redisQueue:  redisQueue,
		sqliteQueue: sqliteQueue,
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

	client, err := grpc.NewClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer client.Close()

	priorities := []int32{int32(task.LowPriority), int32(task.MediumPriority), int32(task.HighPriority)}
	functions := []string{"printHello", "simulateWork"}

	for i := 0; i < 10; i++ {
		priority := priorities[i%len(priorities)]
		functionName := functions[i%len(functions)]
		err := client.SubmitTask(int32(i), functionName, priority)
		if err != nil {
			log.Printf("Failed to submit task: %v", err)
		}
	}

	// Wait for tasks to complete
	time.Sleep(time.Second * 15)
}
