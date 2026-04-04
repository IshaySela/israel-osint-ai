package workerpool

import "sync"

type WorkerPool struct {
	taskQueue   chan func()
	workerCount int
	wg          sync.WaitGroup
}

func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
	return &WorkerPool{
		taskQueue:   make(chan func(), queueSize),
		workerCount: workerCount,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for task := range wp.taskQueue {
				task()
			}
		}()
	}
}

func (wp *WorkerPool) Submit(task func()) { wp.taskQueue <- task }
func (wp *WorkerPool) Wait()              { wp.wg.Wait() }
