package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/ykkalexx/distributed-taskqueue/proto"
	"google.golang.org/grpc"
)

func SubmitTask(addr string, id int32, functionName string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewTaskServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SubmitTask(ctx, &proto.TaskRequest{Id: id, FunctionName: functionName})
	if err != nil {
		return fmt.Errorf("could not submit task: %v", err)
	}

	fmt.Printf("Response: %s\n", r.Message)
	return nil
}