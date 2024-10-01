package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/ykkalexx/distributed-taskqueue/proto"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	client proto.TaskServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}

	c := proto.NewTaskServiceClient(conn)

	return &Client{
		conn: conn,
		client: c,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SubmitTask(id int32, functionName string, priority int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.client.SubmitTask(ctx, &proto.TaskRequest{Id: id, FunctionName: functionName, Priority: priority})
	if err != nil {
		return fmt.Errorf("could not submit task: %v", err)
	}

	fmt.Printf("Response: %s\n", r.Message)
	return nil
}