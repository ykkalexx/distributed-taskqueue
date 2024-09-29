package worker

import (
	"fmt"
	"time"

	"github.com/ykkalexx/distributed-taskqueue/internal/task"
)

func Start(id int, queue *task.Queue) {
	for {
		task, ok := queue.GetTask()
		if !ok {
			time.Sleep(time.Second)
			continue
		}
		fmt.Printf("Worker %d executing task %d\n", id, task.ID)
		err := task.Function()
		if err != nil {
			fmt.Printf("Error executing task %d: %v\n", task.ID, err)
		}
	}
}