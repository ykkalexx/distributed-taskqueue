package loadbalancer

import "sync"

type Worker interface {
	ID() int
}

type LoadBalancer struct {
	workers []Worker
	mutex sync.Mutex
	current int
}

func New() *LoadBalancer {
	return &LoadBalancer{
		workers: make([]Worker, 0),
		current: -1,
	}
}

func (lb *LoadBalancer) AddWorker(worker Worker) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	lb.workers = append(lb.workers, worker)
}

func (lb *LoadBalancer) NextWorker() Worker {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	if len(lb.workers) == 0 {
		return nil
	}

	lb.current = (lb.current + 1) % len(lb.workers)
	return lb.workers[lb.current]
}