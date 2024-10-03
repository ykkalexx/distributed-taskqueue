package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/ykkalexx/distributed-taskqueue/internal/loadbalancer"
	"github.com/ykkalexx/distributed-taskqueue/internal/logger"
	"github.com/ykkalexx/distributed-taskqueue/internal/task"
	"github.com/ykkalexx/distributed-taskqueue/proto"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedTaskServiceServer
	queue task.Queue
	lb    *loadbalancer.LoadBalancer
}

func (s *server) SubmitTask(ctx context.Context, req *proto.TaskRequest) (*proto.TaskResponse, error) {
	t := task.Task{
		ID:           int(req.Id),
		FunctionName: req.FunctionName,
		Priority:     task.Priority(req.Priority),
		MaxRetries:   int(req.MaxRetries),
	}

	err := s.queue.AddTask(t)
	if err != nil {
		logger.Error("Failed to add task: %v", err)
		return &proto.TaskResponse{Success: false, Message: fmt.Sprintf("Failed to add task: %v", err)}, nil
	}

	worker := s.lb.NextWorker()
	if worker != nil {
		logger.Info("Task %d (Priority: %d, Max Retries: %d) assigned to Worker %d", t.ID, t.Priority, t.MaxRetries, worker.ID())
	}

	return &proto.TaskResponse{Success: true, Message: "Task submitted successfully"}, nil
}

func StartServer(queue task.Queue, lb *loadbalancer.LoadBalancer, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterTaskServiceServer(s, &server{queue: queue, lb: lb})

	logger.Info("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
