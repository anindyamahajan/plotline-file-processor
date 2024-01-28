package main

import "sync"

type WorkerPool struct {
	Work    chan []string
	Success bool
	Mutex   sync.Mutex
	WG      sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		Work:    make(chan []string),
		Success: true,
	}

	for i := 0; i < numWorkers; i++ {
		pool.WG.Add(1)
		go worker(pool)
	}

	return pool
}

func (pool *WorkerPool) Reset() {
	pool.Mutex.Lock()
	defer pool.Mutex.Unlock()

	pool.Success = true
	pool.WG = sync.WaitGroup{} // Reset the WaitGroup
}

func worker(pool *WorkerPool) {
	defer pool.WG.Done()
	for data := range pool.Work {
		pool.WG.Add(1) // Increment the counter when starting to process a row
		if !processRow(data) {
			pool.Mutex.Lock()
			pool.Success = false
			pool.Mutex.Unlock()
		}
		pool.WG.Done()
	}
}
