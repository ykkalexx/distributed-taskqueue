package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisQueue struct {
	client *redis.Client
	key    string
}

func NewRedisQueue(addr, password string, db int, key string) (*RedisQueue, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Ping the Redis server to check if the connection is alive
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisQueue{
		client: client,
		key:    key,
	}, nil
}

func (rq *RedisQueue) AddTask(task Task) error {
	ctx := context.Background()
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %v", err)
	}

	// use prioirty as score but subtract retries to lower priority of retried tasks
	score := float64(task.Priority) - float64(task.Retries)*0.1

	// Use priority as score for sorted set
	err = rq.client.ZAdd(ctx, rq.key, &redis.Z{
		Score:  score,
		Member: taskJSON,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to add task to Redis: %v", err)
	}

	return nil
}

func (rq *RedisQueue) GetTask() (Task, bool, error) {
	ctx := context.Background()

	// Use ZPopMax to get and remove the highest scored member (highest priority task)
	result, err := rq.client.ZPopMax(ctx, rq.key).Result()
	if err == redis.Nil || len(result) == 0 {
		return Task{}, false, nil
	} else if err != nil {
		return Task{}, false, fmt.Errorf("failed to get task from Redis: %v", err)
	}

	taskJSON, ok := result[0].Member.(string)
	if !ok {
		return Task{}, false, fmt.Errorf("failed to convert Redis member to string")
	}

	var task Task
	err = json.Unmarshal([]byte(taskJSON), &task)
	if err != nil {
		return Task{}, false, fmt.Errorf("failed to unmarshal task: %v", err)
	}

	return task, true, nil
}

func (rq *RedisQueue) Close() error {
	return rq.client.Close()
}

var _ Queue = (*RedisQueue)(nil)
