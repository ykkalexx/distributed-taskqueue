package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/ykkalexx/distributed-taskqueue/internal/task"
)

func Start(id int, queue interface{}) {
	for {
		var t task.Task
		var ok bool
		var err error

		switch q := queue.(type) {
		case *task.Queue:
			t, ok = q.GetTask()
		case *task.RedisQueue:
			t, ok, err = q.GetTask()
			if err != nil {
				log.Printf("Worker %d error getting task: %v", id, err)
				time.Sleep(time.Second)
				continue
			}
		default:
			log.Printf("Worker %d: unknown queue type", id)
			return
		}

		if !ok {
			time.Sleep(time.Second)
			continue
		}

		fmt.Printf("Worker %d: processing task %d\n", id, t.ID)
		err = t.Function()
		if err != nil {
			fmt.Printf("Error executing task %d: %v\n", t.ID, err)
		}
	}
}