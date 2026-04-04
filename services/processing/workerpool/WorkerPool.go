package workerpool

import (
	"context"
	"sync"

	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
)

type Worker interface {
	Process(ctx context.Context, event models.RawOsintEvent)
}

type WorkerPool struct {
	taskQueue   chan models.RawOsintEvent
	workerCount int
	wg          sync.WaitGroup
}

func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
	return &WorkerPool{
		taskQueue:   make(chan models.RawOsintEvent, queueSize),
		workerCount: workerCount,
	}
}

func (wp *WorkerPool) Start(ctx context.Context, worker Worker) {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for event := range wp.taskQueue {
				worker.Process(ctx, event)
			}
		}()
	}
}

func (wp *WorkerPool) Submit(event models.RawOsintEvent) {
	wp.taskQueue <- event
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
