package workerpool

import (
	"context"
	"fmt"
	"sync"
)

type Task struct {
	Desc string
	Fn   func() error
}

type WorkerPool struct {
	Wg    sync.WaitGroup
	Tasks chan Task
}

func NewWorkerPool(ctx context.Context, concurrency int) *WorkerPool {
	wp := &WorkerPool{
		Tasks: make(chan Task, 1024),
	}

	if concurrency <= 0 {
		concurrency = 1
	}

	wp.Wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		// Starts N workers
		go func(id int) {
			defer wp.Wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case t, ok := <-wp.Tasks:
					// got job to do -> run function
					if !ok {
						// chan is closed, exiting...
						return
					}
					if err := t.Fn(); err != nil {
						fmt.Println("[worker", id, "] send error:", t.Desc, ":", err)
					}
				}
			}
		}(i + 1)
	}
	return wp
}

func (wp *WorkerPool) Submit(t Task) {
	wp.Tasks <- t
}

func (wp *WorkerPool) Close() {
	close(wp.Tasks)
	wp.Wg.Wait()
}
