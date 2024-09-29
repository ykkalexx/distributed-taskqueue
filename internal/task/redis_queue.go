package task

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type RedisQueue struct {
	client *redis.Client
	key string
}

func NewRedisQueue(addr, password string, db int, key string) (*RedisQueue, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    // Ping the Redis server to check if the connection is alive
    if err := client.Ping().Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %v", err)
    }

    return &RedisQueue{
        client: client,
        key:    key,
    }, nil
}

func (rq *RedisQueue) AddTask(task Task) error {
    taskJSON, err := json.Marshal(task)
    if err != nil {
        return fmt.Errorf("failed to marshal task: %v", err)
    }

    err = rq.client.RPush(rq.key, taskJSON).Err()
    if err != nil {
        return fmt.Errorf("failed to add task to Redis: %v", err)
    }

    return nil
}

func (rq *RedisQueue) GetTask() (Task, bool, error) {
    taskJSON, err := rq.client.LPop(rq.key).Result()
    if err == redis.Nil {
        return Task{}, false, nil
    } else if err != nil {
        return Task{}, false, fmt.Errorf("failed to get task from Redis: %v", err)
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